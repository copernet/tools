package main

type dispatchType int

const (
	s2sType dispatchType = 1 << iota
	s2mType
	m2sType
	n2mType		// not to realise at current time
)

func dispatch() {
	spendableCount := len(input)
	if spendableCount == 0 {
		log.Error("There is no spendable transaction.")
	}

	for getDispatchType(m2sType) && !isEmpty() {
		if len(input) < InputLimit {
			break
		}
		m2sTx(true)
	}

	for getDispatchType(s2mType) && !isEmpty() {
		s2mTx(true)
	}

	for getDispatchType(s2sType) && !isEmpty() {
		s2sTx(true)
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
