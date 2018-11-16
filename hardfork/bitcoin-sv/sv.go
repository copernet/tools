package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/bcext/cashutil"
	"github.com/bcext/gcash/chaincfg"
	"github.com/bcext/gcash/chaincfg/chainhash"
	"github.com/bcext/gcash/rpcclient"
	"github.com/bcext/gcash/txscript"
	"github.com/bcext/gcash/wire"
	"github.com/qshuai/tcolor"
	"github.com/shopspring/decimal"
)

var (
	param  = &chaincfg.MainNetParams
	client *rpcclient.Client

	wif            *cashutil.WIF
	sender         cashutil.Address
	feerateDecimal decimal.Decimal
	wait           int

	utxos chan utxo
)

const (
	defaultSignatureSize = 107
	defaultSequence      = 0xffffffff

	// in fact, 540 satoshi is enough.
	defaultP2SHoutputValue = 546

	// to ensure the transaction will be relayed
	addfee = 10
)

type utxo struct {
	hash   *chainhash.Hash
	vout   int
	value  int64
	script []byte
}

func main() {
	for u := range utxos {
		fmt.Println("available utxo:", len(utxos)+1)

		// assemble transaction with necessary elements.
		tx, err := createAssembleTx(u, sender, wif, feerateDecimal)
		if err != nil {
			fmt.Println(tcolor.WithColor(tcolor.Red, "assemble transaction failed: "+err.Error()))
			os.Exit(1)
		}

		fmt.Println("first transaction: ")
		fmt.Println("\thash:", tx.TxHash())
		buf := bytes.NewBuffer(nil)
		err = tx.Serialize(buf)
		if err != nil {
			fmt.Println(tcolor.WithColor(tcolor.Red, "transaction serialize failed: "+err.Error()))
			os.Exit(1)
		}
		fmt.Println("rawtx:", hex.EncodeToString(buf.Bytes()))
		// broadcast the first transaction
		_, err = client.SendRawTransaction(tx, false)
		if err != nil {
			panic(tcolor.WithColor(tcolor.Red, "broadcast transaction failed: "+err.Error()))
		}

		//  =======================
		// shoud add a utxo to channel
		shash := tx.TxHash()
		spendableutxo := utxo{
			hash:   &shash,
			value:  tx.TxOut[1].Value,
			vout:   1,
			script: tx.TxOut[1].PkScript,
		}
		utxos <- spendableutxo

		// -------------------------

		// so now we get a utxo.
		hash := tx.TxHash()
		u := utxo{
			hash:   &hash,
			value:  defaultP2SHoutputValue,
			vout:   0,
			script: tx.TxOut[0].PkScript,
		}
		//  =======================

		spendTx, err := spendAssembleTx(u, wif)
		if err != nil {
			fmt.Println(tcolor.WithColor(tcolor.Red, "assemble transaction failed: "+err.Error()))
			os.Exit(1)
		}
		// broadcast the second transaction
		_, err = client.SendRawTransaction(spendTx, false)
		if err != nil {
			fmt.Println(tcolor.WithColor(tcolor.Red, "broadcast transaction failed: "+err.Error()))
		}

		fmt.Println("second transction:")
		fmt.Println("\thash:", spendTx.TxHash())
		buf2 := bytes.NewBuffer(nil)
		err = tx.Serialize(buf2)
		if err != nil {
			fmt.Println(tcolor.WithColor(tcolor.Red, "transaction serialize failed: "+err.Error()))
			os.Exit(1)
		}
		fmt.Println("\trawtx:", hex.EncodeToString(buf2.Bytes()))

		// leave a blank line
		fmt.Println()

		time.Sleep(time.Duration(wait) * time.Second)
	}
}

// the output with special opcode locate on the first output.
func createAssembleTx(u utxo, sender cashutil.Address,
	wif *cashutil.WIF, feerate decimal.Decimal) (*wire.MsgTx, error) {

	var tx wire.MsgTx
	tx.Version = 1
	tx.LockTime = 0

	// reference a spendable utxo
	outpoint := wire.NewOutPoint(u.hash, uint32(u.vout))
	tx.TxIn = append(tx.TxIn, wire.NewTxIn(outpoint, nil))
	tx.TxIn[0].Sequence = defaultSequence

	tx.TxOut = make([]*wire.TxOut, 2)
	// ========================================================================================
	// p2sh lock script with opcode OP_CHECKDATASIG
	script, err := txscript.NewScriptBuilder().AddOp(txscript.OP_1).
		AddOp(txscript.OP_3).AddOp(txscript.OP_MUL).AddOp(txscript.OP_NUMEQUAL).Script()
	if err != nil {
		return nil, err
	}
	scriptHash := cashutil.Hash160(script)

	pkScript, err := txscript.NewScriptBuilder().AddOp(txscript.OP_HASH160).
		AddData(scriptHash).AddOp(txscript.OP_EQUAL).Script()
	if err != nil {
		return nil, err
	}
	tx.TxOut[0] = &wire.TxOut{PkScript: pkScript, Value: defaultP2SHoutputValue}
	// ========================================================================================

	// add a change output, the offset in output is 1.
	changeScript, err := txscript.PayToAddrScript(sender)
	if err != nil {
		return nil, err
	}
	tx.TxOut[1] = &wire.TxOut{PkScript: changeScript}

	// calculate the chang amount
	txsize := tx.SerializeSize() + defaultSignatureSize
	fee := feerate.Mul(decimal.New(int64(txsize*1e5), 0)).Truncate(0).IntPart()
	outValue := u.value - fee - addfee - defaultP2SHoutputValue
	tx.TxOut[1].Value = outValue

	sourcePkScript, err := txscript.PayToAddrScript(sender)
	if err != nil {
		return nil, err
	}
	// sign the transaction
	return createSign(&tx, []int64{u.value}, sourcePkScript, wif)
}

func createSign(tx *wire.MsgTx, inputValueSlice []int64, pkScript []byte, wif *cashutil.WIF) (*wire.MsgTx, error) {
	for idx, _ := range tx.TxIn {
		sig, err := txscript.RawTxInSignature(tx, idx, pkScript, cashutil.Amount(inputValueSlice[idx]),
			txscript.SigHashAll|txscript.SigHashForkID, wif.PrivKey)
		if err != nil {
			return nil, err
		}
		sig, err = txscript.NewScriptBuilder().AddData(sig).Script()
		if err != nil {
			return nil, err
		}

		pk, err := txscript.NewScriptBuilder().AddData(wif.PrivKey.PubKey().SerializeCompressed()).Script()
		if err != nil {
			return nil, err
		}
		sig = append(sig, pk...)
		tx.TxIn[0].SignatureScript = sig

		// check whether signature is ok or not.
		engine, err := txscript.NewEngine(pkScript, tx, idx, txscript.StandardVerifyFlags,
			nil, nil, inputValueSlice[idx])
		if err != nil {
			return nil, err
		}
		// execution the script in stack
		err = engine.Execute()
		if err != nil {
			return nil, err
		}
	}

	return tx, nil
}

func spendAssembleTx(u utxo, wif *cashutil.WIF) (*wire.MsgTx, error) {
	var tx wire.MsgTx
	tx.Version = 1
	tx.LockTime = 0

	// add a OP_RETURN output
	tx.TxOut = make([]*wire.TxOut, 1)
	script, err := txscript.NewScriptBuilder().AddOp(txscript.OP_RETURN).
		AddData([]byte("the transaction is valid on chain of bitcoin-sv")).
		Script()
	if err != nil {
		return nil, err
	}
	tx.TxOut[0] = wire.NewTxOut(0, script)

	outpoint := wire.NewOutPoint(u.hash, uint32(u.vout))
	tx.TxIn = append(tx.TxIn, wire.NewTxIn(outpoint, nil))
	tx.TxIn[0].Sequence = defaultSequence

	// sign the transaction
	return spendSign(&tx, []int64{u.value}, u.script, wif)
}

func spendSign(tx *wire.MsgTx, inputValueSlice []int64, pkScript []byte, wif *cashutil.WIF) (*wire.MsgTx, error) {
	for idx, _ := range tx.TxIn {
		script, err := txscript.NewScriptBuilder().AddOp(txscript.OP_1).AddOp(txscript.OP_3).AddOp(txscript.OP_MUL).
			AddOp(txscript.OP_NUMEQUAL).Script()
		if err != nil {
			return nil, err
		}
		sig, err := txscript.NewScriptBuilder().AddOp(txscript.OP_3).AddData(script).Script()
		if err != nil {
			return nil, err
		}

		tx.TxIn[idx].SignatureScript = sig

		engine, err := txscript.NewEngine(pkScript, tx, idx, txscript.StandardVerifyFlags|txscript.ScriptEnableCheckDataSig,
			nil, nil, inputValueSlice[idx])
		if err != nil {
			return nil, err
		}

		// verify the signature
		err = engine.Execute()
		if err != nil {
			return nil, err
		}
	}

	return tx, nil
}

func paseUtxo(data []string, owner cashutil.Address) ([]utxo, error) {
	set := make([]utxo, 0, len(data))
	for _, item := range data {
		origin := strings.Split(item, ":")
		if len(origin) != 3 {
			return nil, errors.New("element not enough")
		}

		hash, err := chainhash.NewHashFromStr(origin[0])
		if err != nil {
			return nil, err
		}

		vout, err := strconv.Atoi(origin[1])
		if err != nil {
			return nil, err
		}

		amount, err := strconv.Atoi(origin[2])
		if err != nil {
			return nil, err
		}

		pkscript, err := txscript.PayToAddrScript(owner)
		if err != nil {
			return nil, err
		}

		set = append(set, utxo{
			hash:   hash,
			vout:   vout,
			value:  int64(amount),
			script: pkscript,
		})
	}

	return set, nil
}

func GetRPC(host, user, passwd string) *rpcclient.Client {
	if client != nil {
		return client
	}

	// rpc client instance
	connCfg := &rpcclient.ConnConfig{
		Host:         host,
		User:         user,
		Pass:         passwd,
		HTTPPostMode: true, // Bitcoin core only supports HTTP POST mode
		DisableTLS:   true, // Bitcoin core does not provide TLS by default
	}

	// Notice the notification parameter is nil since notifications are
	// not supported in HTTP POST mode.
	c, err := rpcclient.New(connCfg, nil)
	if err != nil {
		panic(err)
	}

	client = c
	return c
}

func init() {
	privKey := flag.String("privkey", "", "private key of the sender")
	feerate := flag.String("feerate", "0.00001", "the specified feerate for bitcoin cash network")
	second := flag.Int("wait", 60, "the interval for creating transaction")

	host := flag.String("rpchost", "127.0.0.1:8332", "Please input rpc host(ip:port)")
	user := flag.String("rpcuser", "", "Please input your rpc username")
	passwd := flag.String("rpcpassword", "", "Please input your rpc password")
	flag.Parse()

	wait = *second

	if *privKey == "" {
		fmt.Println(tcolor.WithColor(tcolor.Red, "empty private key not allowed"))
		os.Exit(1)
	}

	client = GetRPC(*host, *user, *passwd)

	// parse privkey
	var err error
	wif, err = cashutil.DecodeWIF(*privKey)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "private key format error: "+err.Error()))
		os.Exit(1)
	}

	// get bitcoin address for sender
	pubKey := wif.PrivKey.PubKey()
	pubKeyBytes := pubKey.SerializeCompressed()
	pkHash := cashutil.Hash160(pubKeyBytes)
	sender, err = cashutil.NewAddressPubKeyHash(pkHash, param)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "address encode failed, please check your privkey: "+err.Error()))
		os.Exit(1)
	}

	// parse feerate
	feerateDecimal, err = decimal.NewFromString(*feerate)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "incorrect feerate: "+err.Error()))
		os.Exit(1)
	}

	// init utxo container
	utxos = make(chan utxo, 100)
	// insert some utxos, using hard code,
	// the format:
	// [previous output hash]:[output index]:[value(in satoshi)]
	hashStr := []string{
		//"328040b5b468780eb62d99a1d3da5f1c998ed6d27a08105eadbaaed1b1b98091:0:9996659",
		//"b9b35b7a35ce55193c7bfc201cc4aaf88caa0cf86d36eca2a1f36f90b6a694cc:1:10000000",
		//"5d9c28fc540286c34ba7357398955f6398cc226b3375dd6e241adf5d1f12489a:1:10000000",
		"4cf57ff17b7402d7f08a6860687d12f0d7f55eb8920eb7a3582dcfd431ad6c9d:1:9989860",
		//"2c0dc0d4f3da54b1356c10bf2077ee0875dde63e664c89672a11df395c7d343a:1:9993760",
	}

	ret, err := paseUtxo(hashStr, sender)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "init utxo failed: "+err.Error()))
		os.Exit(1)
	}

	// import utxo element to channel
	for _, item := range ret {
		utxos <- item
	}
}
