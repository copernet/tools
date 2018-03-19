package main

func dispatch() {
	spendableCount := len(input)
	if spendableCount == 0 {
		log.Error("There is no spendable transaction.")
	}

	// create single to single transaction if there are abundant spendable transaction
	// in wallet
	abundantTransactions := conf.DefaultInt("abundant_transactions", AbundantTransactions)
	if spendableCount >= abundantTransactions {
		for !isEmpty() {
			s2sTx(true)
		}
	}

	listUnspentLimit := conf.DefaultInt("exec::list_unspent_limit", DefaultListUnspentLimit)
	iteration :=  listUnspentLimit - spendableCount
	if iteration > 0 {
		log.Info("less input, create more spendable transactions...")
		count := conf.DefaultInt("tx::output_limit", OutputLimit)
		iteration = iteration / count * 2

		for i := 0; i < iteration; i++ {
			s2mTx(true)

			// has too much spendable transactions
			if len(input) > listUnspentLimit*2 {
				break
			}
		}
	}

	if len(lessCoin) > 5000 {
		m2sTx(true)
	}

	for !isEmpty() {
		s2sTx(true)
	}

	// output tip message
	log.Info("Create Transactions Complete!\n")
}

func isEmpty() bool {
	// stop if no input data
	return len(input) == 0
}
