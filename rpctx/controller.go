package main

import "github.com/astaxie/beego/logs"

type dispatchType int

const (
	s2sType dispatchType = 1 << iota
	s2mType
	m2sType
	n2mType
)

func dispatch() {
	// rangeAccount(client)
	inputs(client)

	spendableCount := len(input)
	if spendableCount == 0 {
		logs.Error("There is no spendable transaction.")
	}

	// whether to create transaction recursively
	recursionConf := conf.DefaultBool("recursion", DefaultRecursion)

	if !isEmpty() {
		switch t := getDispatchType(); t {
		case m2sType:
			m2sTx(recursionConf)
		case s2mType:
			s2mTx(recursionConf)
		case s2sType:
			s2sTx(recursionConf)
		case n2mType:
			n2mTx(recursionConf)
		}
	}

	// output tip message
	logs.Info("Create Transactions Complete!\n")
}

func isEmpty() bool {
	// stop if no input data
	return len(input) == 0
}

func getDispatchType() dispatchType {
	dispatch, err := conf.Int("dispatch::type")
	if err != nil {
		panic(err)
	}
	return dispatchType(dispatch)
}
