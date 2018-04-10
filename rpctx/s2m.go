package main

import (
	"math"

	"github.com/astaxie/beego/logs"
	"github.com/btcsuite/btcd/wire"
)

func s2mTx(recursion bool) {
	logs.Info("EXEC s2mTx(%t)", recursion)
	dust := conf.DefaultInt("tx::dust", DefaultDust)
	iteration := conf.DefaultInt("tx::output_limit", OutputLimit)

	for reference, amount := range input {
		// avoid to create a coin with low amount than dust
		// side effect: make coin more nearly amount at the same time
		// todo notice: The number of its tx_out may be not the specified output_limit
		// because of maxSplit judgement
		var maxSplit int
		if dust != 0 {
			maxSplit = int(amount*math.Pow10(8)) / dust
		}
		if maxSplit == 0 {
			continue
		}

		txin := wire.TxIn{
			PreviousOutPoint: wire.OutPoint{
				Hash:  reference.hash,
				Index: reference.index,
			},
			Sequence: 0xffffff, // default value
		}
		s2m.TxIn[0] = &txin

		pkScript := getRandScriptPubKey()
		if pkScript == nil {
			panic("no account in output...")
		}

		bak := iteration
		if maxSplit < bak {
			bak = maxSplit
		}

		splitValue := int(amount*math.Pow10(8))/bak - fee
		if splitValue < 0 {
			continue
		}

		s2m.TxOut = make([]*wire.TxOut, bak)
		for i := 0; i < int(bak); i++ {
			out := wire.TxOut{
				Value:    int64(splitValue), // transaction fee
				PkScript: pkScript,
			}
			s2m.TxOut[i] = &out
		}
		//  no assignment for tx.LockTime(default 0)

		signAndSendTx(s2m, []ref{reference}, int(bak), recursion)
	}
}
