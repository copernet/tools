package main

import (
	"math"

	"github.com/astaxie/beego/logs"
	"github.com/btcsuite/btcd/wire"
)

func m2sTx(recursion bool) {
	logs.Info("EXEC m2sTx(%t)", recursion)

	inputLimit := conf.DefaultInt("tx::input_limit", InputLimit)
	allInOne := conf.DefaultBool("exec::all_in_one", DefaultAllInOne)

	refs := make([]ref, inputLimit)
	counter := 0
	sum := 0.0
	for reference, amount := range input {
		// not enough input items
		if len(input) < inputLimit {
			// all available input will become one item
			if allInOne {
				inputLimit = len(input)
			} else {
				break
			}
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

		give := int(sum*math.Pow10(8)) - inputLimit*fee
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
		m2s.TxIn = m2s.TxIn[:inputLimit]
		signAndSendTx(m2s, refs, 1, recursion)
	}
}
