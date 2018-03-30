package main

import (
	"time"
)

func main() {
	defer client.Shutdown()

	// stop when input is empty or CTRL + C
	for {
		dispatch()

		time.Sleep(time.Duration(interval) * time.Second)
	}
}
