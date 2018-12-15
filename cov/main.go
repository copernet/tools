package main

import (
	"fmt"
	"os"

	"github.com/bcext/cashutil"
	"github.com/bcext/gcash/chaincfg"
)

var paramsMapping = map[string]*chaincfg.Params{
	"mainnet": &chaincfg.MainNetParams,
	"testnet": &chaincfg.TestNet3Params,
	"regtest": &chaincfg.RegressionNetParams,
}

func main() {
	//address := flag.String("address", "", "Please input a bitcoin address string.")
	//targetFormat := flag.String("netparam", "mainnet", "Please input the net identifier: mainnet/testnet/regtest/simnet")
	//
	//flag.Parse()

	args := os.Args
	if len(args) != 3 {
		fmt.Println("param error, please see help information!")
		fmt.Println("Usage: ./cov address netparam[mainnet/testnet/regtest/simnet]")
		return
	}

	var net *chaincfg.Params
	var ok bool
	if net, ok = paramsMapping[args[2]]; !ok {
		fmt.Println("netparam does not exist: please see help information!")
		return
	}

	addr, err := cashutil.DecodeAddress(args[1], net)
	if err != nil {
		fmt.Println("please check address format!")
		return
	}

	// bitcoin cash base32-encoded address:
	// 1. mainnet:  bitcoincash: + 42 bytes = 12 + 42
	// 2. testnet3: bchtest:     + 42 bytes = 8  + 42
	// 3. regtest:  bchreg:      + 42 bytes = 7  + 42
	// if the address does not have the prefix: length = 42
	// if the address have the prefix: max length = 42 + 12
	if len(args[1]) >= 42 && len(args[1]) <= 42+12 {
		fmt.Println(addr.EncodeAddress(false))
	} else {
		fmt.Println(addr.EncodeAddress(true))
	}
}
