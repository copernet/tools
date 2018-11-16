package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/bcext/cashutil"
	"github.com/bcext/gcash/chaincfg"
	"github.com/bcext/gcash/chaincfg/chainhash"
	"github.com/bcext/gcash/txscript"
	"github.com/bcext/gcash/wire"
	"github.com/qshuai/tcolor"
	"github.com/shopspring/decimal"
)

const (
	defaultSignatureSize = 107
	defaultSequence      = 0xffffffff
)

func main() {
	privKey := flag.String("privkey", "", "private key of the sender")
	to := flag.String("to", "", "the bitcoin cash address of receiver")
	hash := flag.String("hash", "", "previous tx hash")
	idx := flag.Int("idx", 0, "previous tx index")
	value := flag.Int("value", 0, "the utxo value")
	pkScript := flag.String("pkscript", "", "the utxo lock script")
	feerate := flag.String("feerate", "0.00001", "the specified feerate for bitcoin cash network")
	flag.Parse()

	if *privKey == "" || *to == "" || *hash == "" || *pkScript == "" || *value == 0 {
		fmt.Println(tcolor.WithColor(tcolor.Red, "arguments are not enough(privkey/to/hash/value required)"))
		os.Exit(1)
	}

	h, err := chainhash.NewHashFromStr(*hash)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "not valid transaction hash: "+err.Error()))
		os.Exit(1)
	}
	// parse privkey
	wif, err := cashutil.DecodeWIF(*privKey)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "privvate key format error: "+err.Error()))
		os.Exit(1)
	}

	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "address encode failed, please check your privkey: "+err.Error()))
		os.Exit(1)
	}

	dst, _ := cashutil.DecodeAddress(*to, &chaincfg.TestNet3Params)

	// parse feerate
	feerateDecimal, err := decimal.NewFromString(*feerate)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "incorrect feerate: "+err.Error()))
		os.Exit(1)
	}

	sourcePkScript, err := hex.DecodeString(*pkScript)
	tx, err := assembleTx(h, int64(*value), uint32(*idx), sourcePkScript, dst, wif, feerateDecimal)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "assemble transaction failed: "+err.Error()))
		os.Exit(1)
	}

	buf := bytes.NewBuffer(nil)
	err = tx.Serialize(buf)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "transaction serialize failed: "+err.Error()))
		os.Exit(1)
	}

	fmt.Println(tcolor.WithColor(tcolor.Green, hex.EncodeToString(buf.Bytes())))
}

func assembleTx(hash *chainhash.Hash, value int64, idx uint32, sourcePkScript []byte, receiver cashutil.Address, wif *cashutil.WIF, feerate decimal.Decimal) (*wire.MsgTx, error) {
	var tx wire.MsgTx
	tx.Version = 1
	tx.LockTime = 0

	tx.TxOut = make([]*wire.TxOut, 1)
	script, err := txscript.PayToAddrScript(receiver)
	if err != nil {
		return nil, err
	}
	tx.TxOut[0] = &wire.TxOut{PkScript: script}

	outpoint := wire.NewOutPoint(hash, idx)
	tx.TxIn = append(tx.TxIn, wire.NewTxIn(outpoint, nil))
	tx.TxIn[0].Sequence = defaultSequence

	txsize := tx.SerializeSize() + defaultSignatureSize

	fee := feerate.Mul(decimal.New(int64(txsize*1e5), 0)).Truncate(0).IntPart()
	outValue := value - fee
	tx.TxOut[0].Value = outValue

	// sign the transaction
	return sign(&tx, []int64{value}, sourcePkScript, wif)
}

func sign(tx *wire.MsgTx, inputValueSlice []int64, pkScript []byte, wif *cashutil.WIF) (*wire.MsgTx, error) {
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

		data := []byte{1, 2}
		hash := sha256.Sum256(data)
		signature, err := wif.PrivKey.Sign(hash[:])
		if err != nil {
			return nil, err
		}
		pk, err := txscript.NewScriptBuilder().AddData(wif.PrivKey.PubKey().SerializeCompressed()).
			AddData(signature.Serialize()).AddData(data).Script()
		if err != nil {
			return nil, err
		}
		sig = append(sig, pk...)
		tx.TxIn[0].SignatureScript = sig

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
