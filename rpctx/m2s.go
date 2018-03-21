package main

import (
	"github.com/btcsuite/btcd/wire"
	"math"
)

func m2sTx(recursion bool) {
	log.Info("EXEC m2sTx(%t)", recursion)

	refs := make([]ref, InputLimit)
	counter := 0
	sum := 0.0
	for reference, amount := range input {
		txin := wire.TxIn{
			PreviousOutPoint: wire.OutPoint{
				Hash:  reference.hash,
				Index: reference.index,
			},
			Sequence: 0xffffff,
		}
		m2s.TxIn[counter] = &txin
		refs[counter] = reference
		counter++
		sum += amount

		if counter < InputLimit {
			continue
		}
		counter = 0

		pkScript := getRandScriptPubKey()
		if pkScript == nil {
			panic("no account in output...")
		}

		give := int64(sum * math.Pow10(8)) - InputLimit * fee
		out := wire.TxOut{
			Value:    give, // transaction fee
			PkScript: pkScript,
		}

		m2s.TxOut[0] = &out

		// reset sum
		sum = 0.0
		signAndSendTx(m2s, refs, 1, recursion)
	}

	s2m.LockTime = 0
}
