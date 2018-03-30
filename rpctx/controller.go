package main

type dispatchType int

const (
	s2sType dispatchType = 1 << iota
	s2mType
	m2sType
	n2mType // not to realise at current time
)

func dispatch() {
	// rangeAccount(client)
	inputs(client)

	spendableCount := len(input)
	if spendableCount == 0 {
		log.Error("There is no spendable transaction.")
	}

	// whether to create transaction recursively
	recursionConf := conf.DefaultBool("recursion", DefaultRecursion)

	if getDispatchType(m2sType) && !isEmpty() {
		m2sTx(recursionConf)
	}

	if getDispatchType(s2mType) && !isEmpty() {
		s2mTx(recursionConf)
	}

	if getDispatchType(s2sType) && !isEmpty() {
		s2sTx(recursionConf)
	}

	if getDispatchType(n2mType) && !isEmpty() {
		n2mTx(recursionConf)
	}

	// output tip message
	log.Info("Create Transactions Complete!\n")
}

func isEmpty() bool {
	// stop if no input data
	return len(input) == 0
}

func getDispatchType(t dispatchType) bool {
	dispatch := must(conf.Int("dispatch::type"))
	return dispatchType(dispatch.(int))&t == t
}
