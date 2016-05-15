package main

import (
	"fmt"
)

func main() {

	// Setup channels for the pipeline
	newTrans := make(chan Transaction, 10)

	// Start the pipeline
	go Link(newTrans)

	for t := range newTrans {
		fmt.Println(t)
	}
}
