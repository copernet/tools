package main

import (
	"math"
	"math/rand"

	"github.com/btcsuite/btcd/wire"
)

func n2mTx(recursion bool) {
	log.Info("EXEC m2sTx(%t)", recursion)

	inputLimit := conf.DefaultInt("input_limit", InputLimit)
	iteration := conf.DefaultInt("output_limit", OutputLimit)

	// plus 1 to insure the result never be zero
	realInputs := rand.Intn(inputLimit) + 1
	refs := make([]ref, realInputs)

	// construct txin start------------
	counter := 0
	sum := 0.0
	for reference, amount := range input {
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
		n2m.TxIn = append(n2m.TxIn, &txin)
		refs[counter] = reference
		counter++
		sum += amount

		if counter < realInputs {
			continue
		}
		counter = 0

		// construct txin end--------------
		pkScript := getRandScriptPubKey()
		if pkScript == nil {
			panic("no account in output...")
		}

		splitValue := int(amount*math.Pow10(8))/iteration - fee
		if splitValue < 0 {
			continue
		}

		s2m.TxOut = make([]*wire.TxOut, iteration)
		for i := 0; i < int(iteration); i++ {
			out := wire.TxOut{
				Value:    int64(splitValue), // transaction fee
				PkScript: pkScript,
			}
			s2m.TxOut[i] = &out
		}
		//  no assignment for tx.LockTime(default 0)

		// reset sum
		sum = 0.0
		signAndSendTx(m2s, refs, 1, recursion)
	}
}
