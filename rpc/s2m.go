package main

import (
	"github.com/btcsuite/btcd/wire"
	"math"
)

func s2mTx(recursion bool) {
	log.Info("EXEC s2mTx(%t)", recursion)

	for reference, amount := range input {
		txin := wire.TxIn{
			PreviousOutPoint: wire.OutPoint{
				Hash:  reference.hash,
				Index: reference.index,
			},
			Sequence: 0xffffff,
		}
		s2m.TxIn[0] = &txin

		pkScript := getRandScriptPubKey()
		if pkScript == nil {
			panic("no account in output...")
		}

		dust := conf.DefaultInt("tx::dust", DefaultDust)
		maxSplit := int(amount * math.Pow10(8)) / dust

		var iteration int64 = OutputLimit
		if maxSplit < OutputLimit {
			iteration = int64(maxSplit)
		}

		splitValue := int64(amount * math.Pow10(8)) / iteration - iteration * fee
		for i := 0; i < OutputLimit; i++ {
			out := wire.TxOut{
				Value:    splitValue, // transaction fee
				PkScript: pkScript,
			}

			s2m.TxOut = append(s2m.TxOut, &out)
		}

		signAndSendTx(s2m, []ref{reference}, int(iteration), recursion)
	}

	s2m.LockTime = 0
}
