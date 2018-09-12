package main

import (
	"BarkChain/blockchain"
	"fmt"

)

func main()  {

	bc := blockchain.NewBlockChain()
	bc.NewBlock("Send 4 USD to Fan")
	bc.NewBlock("send 5 USD to Fan")

	// view
	//for index, b := range bc.Blocks {
	//	fmt.Printf("%d, %x , %x,  %s\n",index, string(b.BlockHeader.PrevHash), string(b.BlockHeader.CurHash), b.Data)
	//}

	bci := bc.Iterator()
	for block := bci.Next(); block != nil ; block = bci.Next() {
		fmt.Printf(" %x , %x,  %s\n",string(block.BlockHeader.PrevHash), string(block.BlockHeader.CurHash), block.Data)
	}

	// validate
	bci = bc.Iterator()
	for block := bci.Next(); block != nil ; block = bci.Next() {
		proof := blockchain.NewProofOfWork(block, blockchain.POW_TARGET)
		fmt.Printf("%x , %x,  %s,  %v\n",string(block.BlockHeader.PrevHash), string(block.BlockHeader.CurHash), block.Data, proof.Validate())
	}

	bc.Db.Close()
}
