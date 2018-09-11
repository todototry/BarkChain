package main

import (
	"BarkChain/blockchain"
	"fmt"
)

func main()  {
	bc := blockchain.NewBlockChain()
	bc.NewBlock("send 3 Yuan to Eva")
	bc.NewBlock("send 2 Yuan to Fan")

	// view
	for index, b := range bc.Blocks {
		fmt.Printf("%d, %x , %x,  %s\n",index, string(b.BlockHeader.PrevHash), string(b.BlockHeader.CurHash), b.Data)
	}

	// validate
	for index, b := range bc.Blocks {
		proof := blockchain.NewProofOfWork(b, blockchain.POW_TARGET)
		fmt.Printf("%d, %x , %x,  %s  , %v \n",index, string(b.BlockHeader.PrevHash), string(b.BlockHeader.CurHash), b.Data, proof.Validate() )
	}


}
