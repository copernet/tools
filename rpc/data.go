package main

import (
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/rpcclient"
	"encoding/hex"
	"math"
)

// inputs() function will stop this program via panic exception
// because origin spendable tx will be empty if any error occur.
func inputs(client *rpcclient.Client) {
	log.Info("starting acquire data...")

	dust, err := conf.Int("tx::dust")
	if err != nil {
		dust = DefaultDust
	}

	limitCoin, err := conf.Int("tx::limit_coin")
	if err != nil {
		dust = DefaultDust
	}

	// rpc requests to get unspent coin list
	lu, err := client.ListUnspent()
	if err != nil {
		panic(err)
	}

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

		limitCoinConvert := float64(limitCoin) * math.Pow10(-8.0)
		if item.Amount > dustConvert {

			input[r] = item.Amount

			scriptPubKey, _ := hex.DecodeString(item.ScriptPubKey)
			if err != nil {
				panic(err)
			}
			output[item.Address] = scriptPubKey
		} else if item.Amount > limitCoinConvert{
			lessCoin[r] = item.Amount
		}
	}
	log.Info("input: %d, output: %d, lessCoin: %d", len(input), len(output), len(lessCoin))
}
