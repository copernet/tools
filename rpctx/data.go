package main

import (
	"encoding/hex"
	"math"

	"github.com/astaxie/beego/logs"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
)

// inputs() function will stop this program via panic exception
// because origin spendable tx will be empty if any error occur.
func inputs(client *rpcclient.Client) {
	logs.Info("starting acquire data...")

	dust := conf.DefaultInt("tx::dust", DefaultDust)

	// rpc requests to get unspent coin list
	lu, err := client.ListUnspent()
	if err != nil {
		panic(err)
	}

	var lessCoin int
	for _, item := range lu {
		if item.Vout >= 255 {
			break
		}

		// skip if the balance of this bitcoin address is too low.
		// bitcoin client will shoots out error: "insufficient priority"
		// but it works if set settxfee = 0 via rpc command, like this:
		// bitcoin-cli settxfee 0
		hash, err := chainhash.NewHashFromStr(item.TxID)
		if err != nil {
			continue
		}

		r := ref{
			hash:  *hash,
			index: item.Vout,
		}

		// convert Satoshi to BCH
		dustConvert := float64(dust) * math.Pow10(-8.0)

		if item.Amount > dustConvert {

			input[r] = item.Amount

			scriptPubKey, _ := hex.DecodeString(item.ScriptPubKey)
			if err != nil {
				panic(err)
			}
			output[item.Address] = scriptPubKey
		} else {
			lessCoin++
			if getDispatchType() == m2sType {
				// add the ref transaction only in m2sType
				input[r] = item.Amount
			}
		}
	}
	logs.Info("input: %d, output: %d, lessCoin: %d", len(input), len(output), lessCoin)
}
