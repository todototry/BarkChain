package main

import (
	"BarkChain/blockchain"
	"BarkChain/cmd"
)

func main()  {

	bc := blockchain.NewBlockChain()
	defer bc.Db.Close()

	cli := cmd.NewCli(bc)
	cli.Run()

}
