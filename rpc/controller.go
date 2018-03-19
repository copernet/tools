package main

import (
	"fmt"
)

func dispatch() {
	spendableCount := len(input)
	if spendableCount == 0 {
		log.Error("There is no spendable transaction.")
	}

	listunspentLimit := conf.DefaultInt("exec::listunspent_limit", DefaultListunspentLimit)
	iteration :=  listunspentLimit - spendableCount
	if iteration > 0 {
		log.Info("less input, create more spendable transactions...")
		count := conf.DefaultInt("tx::output_limit", OutputLimit)
		iteration = iteration / count * 2

		for i := 0; i < iteration; i++ {
			s2mTx(true)

			// has too much spendable transactions
			if len(input) > listunspentLimit*2 {
				break
			}
		}
	}

	if len(lessCoin) > 5000 {
		m2sTx(true)
	}

	for {
		s2sTx(false)

		// stop if no input data
		if len(input) == 0 {
			break
		}
	}

	// output tip message
	fmt.Println("Create Transactions Complate!\n")
}
