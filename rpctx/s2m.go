package main

import (
	"github.com/btcsuite/btcd/wire"
	"math"
)

func s2mTx(recursion bool) {
	log.Info("EXEC s2mTx(%t)", recursion)
	dust := conf.DefaultInt("tx::dust", DefaultDust)

	for reference, amount := range input {
		txin := wire.TxIn{
			PreviousOutPoint: wire.OutPoint{
				Hash:  reference.hash,
				Index: reference.index,
			},
			Sequence: 0xffffff,		// default value
		}
		s2m.TxIn[0] = &txin

		pkScript := getRandScriptPubKey()
		if pkScript == nil {
			panic("no account in output...")
		}

		// avoid to create a coin with low amount than dust
		maxSplit := int(amount * math.Pow10(8)) / dust

		if maxSplit == 0 {
			continue
		}

		var iteration int64 = OutputLimit
		if maxSplit < OutputLimit {
			iteration = int64(maxSplit)
		}

		splitValue := int64(amount * math.Pow10(8)) / iteration - fee
		if splitValue < 0 {
			continue
		}

		s2m.TxOut = make([]*wire.TxOut, iteration)
		for  i := 0; i < int(iteration); i++ {
			out := wire.TxOut{
				Value:    splitValue, // transaction fee
				PkScript: pkScript,
			}
			s2m.TxOut[i] = &out
		}
		//! no assignment for tx.LockTime(default 0)

		signAndSendTx(s2m, []ref{reference}, int(iteration), recursion)
	}
}
