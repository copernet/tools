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
	lessCoinValue := conf.DefaultInt("less_coin_limit", LessCoinLimit)
	inputLimit := conf.DefaultInt("input_limit",InputLimit)
	for reference, amount := range input {
		// aggregate many less coins in one other than abundant coin item
		if amount * math.Pow10(8) > float64(lessCoinValue) {
			continue
		}

		// not enough input items
		if len(input) < inputLimit {
			break
		}

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

		if counter < inputLimit {
			continue
		}
		counter = 0

		give := int64(sum * math.Pow10(8)) - InputLimit * fee
		if give <0 {
			continue
		}
		pkScript := getRandScriptPubKey()
		if pkScript == nil {
			panic("no account in output...")
		}
		out := wire.TxOut{
			Value:    give, // transaction fee
			PkScript: pkScript,
		}
		m2s.TxOut[0] = &out
		//! no assignment for tx.LockTime(default 0)

		// reset sum
		sum = 0.0
		signAndSendTx(m2s, refs, 1, recursion)
	}

	s2m.LockTime = 0
}
