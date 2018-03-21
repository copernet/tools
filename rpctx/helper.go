package main

func must(value interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}

	return value
}
