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

	originInputLimit := conf.DefaultInt("tx::input_limit", InputLimit)
	originIteration := conf.DefaultInt("tx::output_limit", OutputLimit)

	refs := make([]ref, 0, originInputLimit)

	// construct txin start------------
	counter := 0
	sum := 0.0
	for reference, amount := range input {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		inputLimit := r.Intn(originInputLimit)
		iteration := r.Intn(originIteration)

		if inputLimit == 0 {
			inputLimit = 1
		}
		if iteration == 0 {
			iteration = 1
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
		n2m.TxIn = append(n2m.TxIn, &txin)
		refs = append(refs, reference)

		counter++
		sum += amount

		if counter < inputLimit {
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
			n2m.TxIn = n2m.TxIn[:0] // clean TxIn element
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
		refs = refs[:0]
	}
}
