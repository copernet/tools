package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/bcext/cashutil"
	"github.com/bcext/gcash/chaincfg"
	"github.com/bcext/gcash/chaincfg/chainhash"
	"github.com/bcext/gcash/txscript"
	"github.com/bcext/gcash/wire"
	"github.com/qshuai/tcolor"
	"github.com/tidwall/gjson"
)

const (
	// address pair
	bech32Address = "bchtest:qqwpvaha3leydercn7kckkuh9ufxaplcmsn48e8v3m"
	base58Address = "mi5U8JnLMLiVrms3mW9YNvz5nSYC57Q7G9"

	privkey = "cRL6HJZYSF1JMUSyP6PsKMRD9PhCS1acUSoKWh9p5Bf5iY4SPq5j"

	// 10 satoshi/byte
	feeRate = 10

	defaultSignatureSize = 107

	defaultSequence = 0xffffffff
)

func main() {
	args := os.Args
	if len(args) <= 1 {
		fmt.Println(tcolor.WithColor(tcolor.Red, "Please a message needed to send"))
		os.Exit(1)
	}

	balance, err := getBalance(base58Address)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "Sorry, get balance of the specified address failed"))
		os.Exit(1)
	}

	// unnecessary to create a transaction if balance is lower than 0.001 BCH
	if balance < 100000 {
		fmt.Println(tcolor.WithColor(tcolor.Red, "Sorry, your balance is insufficient"))
		os.Exit(1)
	}

	utxo, err := getUnspent(base58Address, 1)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "Sorry, get utxo of the specified failed"))
		os.Exit(1)
	}

	scriptHash, err := getPkScriptHash(bech32Address)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "Please check your bech32 address"))
		os.Exit(1)
	}
	pkScript := getP2pkhScript(scriptHash)

	msg := args[1]
	msgScript := txscript.NewScriptBuilder()
	msgScript.AddOp(txscript.OP_RETURN).AddData([]byte(msg))
	msgBytes, err := msgScript.Script()
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "Make OP_RETURN script error:"+err.Error()))
		os.Exit(1)
	}

	// parse privkey
	wif, err := cashutil.DecodeWIF(privkey)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "Privkey format error"))
		os.Exit(1)
	}

	tx, err := assembleTx(utxo, msgBytes, pkScript, wif)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "Assemble transaction or sign error:"+err.Error()))
		os.Exit(1)
	}

	buf := bytes.NewBuffer(nil)
	err = tx.Serialize(buf)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "Transaction serialize error:"+err.Error()))
		os.Exit(1)
	}
	// output result
	fmt.Println("txhash:         ", tcolor.WithColor(tcolor.Green, tx.TxHash().String()))
	fmt.Println("raw transaction:", tcolor.WithColor(tcolor.Green, hex.EncodeToString(buf.Bytes())))
}

func getPkScriptHash(address string) ([]byte, error) {
	addr, err := cashutil.DecodeAddress(address, &chaincfg.TestNet3Params)
	if err != nil {
		return nil, err
	}

	return addr.ScriptAddress(), nil
}

func getP2pkhScript(scriptHash []byte) []byte {
	pkScript := txscript.NewScriptBuilder().AddOp(txscript.OP_DUP).AddOp(txscript.OP_HASH160).
		AddData(scriptHash).AddOp(txscript.OP_EQUALVERIFY).AddOp(txscript.OP_CHECKSIG)

	// ignore error because the specified address is checked
	bs, _ := pkScript.Script()

	return bs
}

func assembleTx(utxo string, msgBytes []byte, pkScript []byte, wif *cashutil.WIF) (*wire.MsgTx, error) {
	var tx wire.MsgTx
	tx.Version = 1
	tx.LockTime = 0

	tx.TxOut = make([]*wire.TxOut, 2)
	tx.TxOut[0] = &wire.TxOut{PkScript: pkScript}
	tx.TxOut[1] = &wire.TxOut{PkScript: msgBytes, Value: 0}

	unspentList := gjson.Get(utxo, "data.list").Array()
	var inputValue float64
	var inputValueSlice []int64
	for i := 0; i < len(unspentList); i++ {
		// the coin value
		value := unspentList[i].Get("value").Float()
		if value <= 0 {
			continue
		}
		inputValue += value
		inputValueSlice = append(inputValueSlice, int64(value))

		hashStr := unspentList[i].Get("tx_hash").String()
		// ignore error because of trusting the API result
		hash, _ := chainhash.NewHashFromStr(hashStr)
		index := unspentList[i].Get("tx_output_n").Int()

		txIn := wire.TxIn{
			PreviousOutPoint: *wire.NewOutPoint(hash, uint32(index)),
			Sequence:         defaultSequence,
		}
		tx.TxIn = append(tx.TxIn, &txIn)

		actualFeeRate := value / float64(tx.SerializeSize()+defaultSignatureSize*(i+1))
		if actualFeeRate < feeRate {
			continue
		}

		fee := (tx.SerializeSize() + defaultSignatureSize*(i+1)) * feeRate
		redeemAmount := int(inputValue) - fee
		tx.TxOut[0].Value = int64(redeemAmount)

		break
	}

	// sign the transaction
	return sign(&tx, inputValueSlice, pkScript, wif)
}
