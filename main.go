package main

import (
	"BarkChain/blockchain"
	"fmt"
)

func main()  {
	bc := blockchain.NewBlockChain()
	bc.NewBlock("send 3 Yuan to Eva")
	bc.NewBlock("send 2 Yuan to Fan")

	for index, b := range bc.Blocks {
		fmt.Println(index, b)
	}

}
