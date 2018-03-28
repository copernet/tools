package main

import (
	"math"
	"math/rand"

	"github.com/btcsuite/btcd/wire"
)

func m2sTx(recursion bool) {
	log.Info("EXEC m2sTx(%t)", recursion)

	inputLimit := conf.DefaultInt("input_limit", InputLimit)

	// plus 1 to insure the result never be zero
	realInputs := rand.Intn(inputLimit) + 1
	refs := make([]ref, realInputs)
	counter := 0
	sum := 0.0
	lessCoinValue := conf.DefaultInt("less_coin_limit", LessCoinLimit)
	for reference, amount := range input {
		// aggregate many less coins in one other than abundant coin item
		if amount*math.Pow10(8) < float64(lessCoinValue) {
			continue
		}

		// not enough input items
		if len(input) < realInputs {
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

		if counter < realInputs {
			continue
		}
		counter = 0

		give := int(sum*math.Pow10(8)) - realInputs*fee
		if give < 0 {
			continue
		}
		pkScript := getRandScriptPubKey()
		if pkScript == nil {
			panic("no account in output...")
		}
		out := wire.TxOut{
			Value:    int64(give), // transaction fee
			PkScript: pkScript,
		}
		m2s.TxOut[0] = &out
		// no assignment for tx.LockTime(default 0)

		// reset sum
		sum = 0.0
		signAndSendTx(m2s, refs, 1, recursion)
	}
}
