package main

import (
	"github.com/btcsuite/btcd/wire"
	"math"
)

func m2sTx(recursion bool) {
	log.Info("EXEC m2sTx(%t)", recursion)

	refs := make([]ref, 0, InputLimit)
	for reference, amount := range input {
		counter := 0
		sum := 0.0

		txin := wire.TxIn{
			PreviousOutPoint: wire.OutPoint{
				Hash:  reference.hash,
				Index: reference.index,
			},
			Sequence: 0xffffff,
		}
		m2s.TxIn = append(m2s.TxIn, &txin)
		counter++
		sum += amount

		refs = append(refs, reference)
		if counter < InputLimit {
			continue
		}


		pkScript := getRandScriptPubKey()
		if pkScript == nil {
			panic("no account in output...")
		}

		give := int64(sum * math.Pow10(8)) - InputLimit * fee
		out := wire.TxOut{
			Value:    give, // transaction fee
			PkScript: pkScript,
		}

		s2s.TxOut[0] = &out

		signAndSendTx(s2m, refs, 1, recursion)
	}

	s2m.LockTime = 0
}
