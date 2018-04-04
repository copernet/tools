package main

import (
	"math"

	"github.com/astaxie/beego/logs"
	"github.com/btcsuite/btcd/wire"
)

func s2sTx(recursion bool) {
	logs.Info("EXEC s2sTx(%t)", recursion)

	for reference, amount := range input {
		give := int(amount*math.Pow10(8)) - fee
		// Discard this transaction if out value less than zero, so that
		// fee rate is nearly equal to each other for pack into a block
		// easily!
		// if creating a transaction out value less than zero, the client
		// will throw an error "-26: 16: bad-txns-vout-negative"
		if give < 0 { // Top Priority exec for optimizing
			continue
		}

		pkScript := getRandScriptPubKey()
		if pkScript == nil {
			panic("no account in output...")
		}
		out := wire.TxOut{
			Value:    int64(give),
			PkScript: pkScript,
		}
		s2s.TxOut[0] = &out

		txin := wire.TxIn{
			PreviousOutPoint: wire.OutPoint{
				Hash:  reference.hash,
				Index: reference.index,
			},
			Sequence: 0xffffff,
		}
		s2s.TxIn[0] = &txin
		// no assignment for tx.LockTime(default 0)

		signAndSendTx(s2s, []ref{reference}, 1, recursion)
	}
}
