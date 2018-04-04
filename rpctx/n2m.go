package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/btcsuite/btcd/wire"
)

func n2mTx(recursion bool) {
	logs.Info("EXEC n2mTx(%t)", recursion)

	inputLimit := conf.DefaultInt("tx::input_limit", InputLimit)
	iteration := conf.DefaultInt("tx::output_limit", OutputLimit)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// plus 1 to insure the result never be zero
	realInputs := r.Intn(inputLimit) + 1
	iteration = r.Intn(iteration) + 1
	refs := make([]ref, inputLimit)

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
		if splitValue <= 0 {
			continue
		}

		for i := 0; i < int(iteration); i++ {
			out := wire.TxOut{
				Value:    int64(splitValue), // transaction fee
				PkScript: pkScript,
			}
			n2m.TxOut = append(n2m.TxOut, &out)
		}
		// no assignment for tx.LockTime(default 0)

		// reset sum
		sum = 0.0
		signAndSendTx(n2m, refs, iteration, recursion)

		// reuse memory space partly
		n2m.TxIn = n2m.TxIn[:0]
		n2m.TxOut = n2m.TxOut[:0]

		// plus 1 to insure the result never be zero
		realInputs = r.Intn(inputLimit) + 1
		iteration = r.Intn(iteration) + 1
	}
}
