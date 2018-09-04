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
		fmt.Printf("%d, %x , %x,  %s\n",index, string(b.BlockHeader.PrevHash), string(b.BlockHeader.CurHash), b.Data)
	}

}
