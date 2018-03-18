package main

import (
	"github.com/btcsuite/btcd/wire"
	"math"
)

func s2sTx(recursion bool) {
	log.Info("EXEC s2sTx(%t)", recursion)

	for reference, amount := range input {
		txin := wire.TxIn{
			PreviousOutPoint: wire.OutPoint{
				Hash:  reference.hash,
				Index: reference.index,
			},
			Sequence: 0xffffff,
		}
		s2s.TxIn[0] = &txin

		pkScript := getRandScriptPubKey()
		if pkScript == nil {
			panic("no account in output...")
		}

		give := int64(amount*math.Pow10(8)) - fee
		out := wire.TxOut{
			Value:    give, // transaction fee
			PkScript: pkScript,
		}

		s2s.TxOut[0] = &out

		signAndSendTx(s2s, []ref{reference}, 1, recursion)
	}
}
