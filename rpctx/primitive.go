package main

import (
	"time"

	"github.com/astaxie/beego/config"
	"github.com/astaxie/beego/logs"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil/base58"
)

const (
	DefaultDust      = 1000
	DefaultLimitCoin = 10

	InputLimit  = 50
	OutputLimit = 50

	DefaultFee = 0

	DefaultListUnspentLimit = 10000
	AbundantTransactions    = 60000
	LessCoinLimit           = 5000
	DefaultRecursion        = true
	DefaultAllInOne         = true

	DefaultInterval   = 600 // second
	DefaultTxInterval = 100 // millisecond
)

// global variables
var (
	conf config.Configer

	// store available input and output
	input  = make(coin)
	output = make(map[string][]byte)
	fee    int

	client *rpcclient.Client
	// successful transaction
	count = 0

	// rpc listunspent interval
	interval = 0
	// send transaction interval
	txinterval = 0

	s2s *wire.MsgTx
	s2m *wire.MsgTx
	m2s *wire.MsgTx
	n2m *wire.MsgTx
)

type ref struct {
	hash  chainhash.Hash
	index uint32
}

type coin map[ref]float64

func init() {
	logs.Info("app init start...")
	// configuration setting
	var err error
	conf, err = config.NewConfig("ini", "conf/app.conf")
	if err != nil {
		panic(err)
	}

	// log setting
	logs.SetLogger("console")
	// log.SetLogger(logs.AdapterFile, `{"filename":"log/btcrpc.log"}`)
	// if must(conf.Bool("log::async")).(bool) {
	// 	log.Async(1e3)
	// }

	// get transaction fee from configuration
	fee = conf.DefaultInt("tx::fee", DefaultFee)
	interval = conf.DefaultInt("dispatch::interval", DefaultInterval)
	txinterval = conf.DefaultInt("dispatch::txinterval", DefaultTxInterval)

	client = Client()

	{
		// object reuse transaction memory space
		s2s = wire.NewMsgTx(1)
		s2s.TxIn = make([]*wire.TxIn, 1)
		s2s.TxOut = make([]*wire.TxOut, 1)

		s2m = wire.NewMsgTx(1)
		s2m.TxIn = make([]*wire.TxIn, 1)

		m2s = wire.NewMsgTx(1)
		inputLimit := conf.DefaultInt("tx::input_limit", InputLimit)
		m2s.TxIn = make([]*wire.TxIn, inputLimit)
		m2s.TxOut = make([]*wire.TxOut, 1)

		n2m = wire.NewMsgTx(1)
		n2m.TxIn = make([]*wire.TxIn, 0)
		n2m.TxOut = make([]*wire.TxOut, 0)
	}

	logs.Info("app init complete!")
}

func signAndSendTx(msg *wire.MsgTx, refs []ref, outs int, recursion bool) {
	// rpc requests signing a raw transaction and gets returned signed transaction,
	// or get null and a err reason
	time.Sleep(time.Duration(txinterval) * time.Millisecond)
	signedTx, _, err := client.SignRawTransaction(msg)
	if err != nil {
		logs.Error(err.Error())
	}

	// todo sign tx in app no to bother client rpc(optimize)
	// btc transaction signature algorithm is different from bch, so
	// following code is invalid.
	// rawPriv, _ := hex.DecodeString("**************")
	// prikey,_ := btcec.PrivKeyFromBytes(btcec.S256(), rawPriv)
	// for idx, _ := range msg.TxIn{
	// 	b, err := txscript.SignatureScript(msg,idx,msg.TxOut[0].PkScript,65503,prikey,true)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	msg.TxIn[idx].SignatureScript = b
	// }

	// rpc request send a signed transaction, it will return a error if there are any
	// error
	txhash, err := client.SendRawTransaction(signedTx, true)

	// remove transactions that have been consumed whether success or not
	removeInputRecursion(refs)
	if err != nil {
		logs.Error(err.Error())
	} else {
		// recursion tx
		if recursion {
			reference := ref{}
			reference.hash = *txhash
			for i := 0; i < outs; i++ {
				reference.index = uint32(i)
				input[reference] = float64(msg.TxOut[i].Value) * 1e-8
			}
		}

		count++
		logs.Info("Create a transaction success, NO.%d, txhash: %s", count, txhash.String())
	}
}

// map return random item
func getRandScriptPubKey() []byte {
	for _, item := range output {
		return item
	}
	return nil
}

func removeInputRecursion(refs []ref) {
	for _, item := range refs {
		delete(input, item)
	}
}

// todo spent to different addresses, support addresses with known ScriptPubKey
func rangeAccount(client *rpcclient.Client) {
	addresses, err := client.GetAddressesByAccount("")
	if err != nil {
		panic(err)
	}

	for _, item := range addresses {
		ret, _, err := base58.CheckDecode(item.String())
		if err != nil {
			panic(err)
		}

		final, err := txscript.NewScriptBuilder().AddOp(txscript.OP_DUP).AddOp(txscript.OP_HASH160).
			AddData(ret).AddOp(txscript.OP_EQUALVERIFY).AddOp(txscript.OP_CHECKSIG).
			Script()

		if err != nil {
			panic(err)
		}
		output[item.String()] = final
	}
}
