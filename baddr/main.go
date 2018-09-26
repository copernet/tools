package main

import (
	"fmt"
	"os"

	"github.com/bcext/cashutil"
	"github.com/bcext/gcash/btcec"
	"github.com/bcext/gcash/chaincfg"
	"github.com/qshuai/tcolor"
)

var net = map[string]*chaincfg.Params{
	"mainnet": &chaincfg.MainNetParams,
	"testnet": &chaincfg.TestNet3Params,
	"regtest": &chaincfg.RegressionNetParams,
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println(tcolor.WithColor(tcolor.Red, "Usage: addr mainnet/testnet/regtest"))
		return
	}

	var n *chaincfg.Params
	n, ok := net[os.Args[1]]
	if !ok {
		fmt.Println(tcolor.WithColor(tcolor.Red, os.Args[1]+" not existed, should select from mainnet/testnet/regtest"))
		return
	}

	priv, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "Create private key failed: "+err.Error()))
		return
	}

	wif, err := cashutil.NewWIF(priv, n, true)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "Convert wif format private key failed: "+err.Error()))
		return
	}

	pubKey := priv.PubKey()
	pubKeyHash := cashutil.Hash160(pubKey.SerializeCompressed())
	addr, err := cashutil.NewAddressPubKeyHash(pubKeyHash, n)
	if err != nil {
		fmt.Println(tcolor.WithColor(tcolor.Red, "Generate a new bitcoin-cash address failed: "+err.Error()))
		return
	}

	fmt.Println("privkey key:           ", tcolor.WithColor(tcolor.Green, wif.String()))
	fmt.Println("base58 encoded address:", tcolor.WithColor(tcolor.Green, addr.EncodeAddress(false)))
	fmt.Println("bech32 encoded address:", tcolor.WithColor(tcolor.Green, addr.EncodeAddress(true)))
}
