package cmd

import (
	"BarkChain/blockchain"
	"fmt"
	"flag"
	"os"
)

const (
	CMD_ADD_BLOCK = "addblock"
	CMD_ADD_BLOCK_OPTION = "data"
	CMD_PRINT_BLOCKCHAIN = "printchain"
)

type Cli struct {
	blockchain *blockchain.BlockChain
}

func NewCli(bc *blockchain.BlockChain) *Cli {
	return &Cli{bc}
}



func (cli *Cli) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  addblock -data BLOCK_DATA - add a block to the blockchain")
	fmt.Println("  printchain - print all the blocks of the blockchain")
}

func (cli *Cli) validateArgs() {
	if len(os.Args) < 2 {
		cli.printUsage()
		os.Exit(1)
	}
}


func (cli *Cli) Run() {
	cli.validateArgs()

	addBlockSubCMD := flag.NewFlagSet(CMD_ADD_BLOCK, flag.ExitOnError)
	printChainSubCMD := flag.NewFlagSet(CMD_PRINT_BLOCKCHAIN, flag.ExitOnError)

	// define param for addblock option
	var data string
	addBlockSubCMD.StringVar(&data, CMD_ADD_BLOCK_OPTION, "", "add a new block with content of data")

	switch os.Args[1] {
	case CMD_ADD_BLOCK:
		addBlockSubCMD.Parse(os.Args[2:])
	case CMD_PRINT_BLOCKCHAIN:
		printChainSubCMD.Parse(os.Args[2:])
	default:
		cli.printUsage()
	}

	if addBlockSubCMD.Parsed() {
		if data != "" {
			cli.AddBlock(data)
		}
	}

	if printChainSubCMD.Parsed() {
		cli.PrintBlockChain()
	}

}

func (cli *Cli) AddBlock(data string) {
	cli.blockchain.NewBlock(data)

	// print out.
	newBlock, _ := cli.blockchain.FindBlock(cli.blockchain.Tip)
	fmt.Printf("%v \n", newBlock)
}

func (cli *Cli) PrintBlockChain() {

	// print out all.
	bci := cli.blockchain.Iterator()
	for block := bci.Next(); block != nil ; block=bci.Next() {
		fmt.Printf(" %x , %x,  %s, %b \n",string(block.BlockHeader.PrevHash), string(block.BlockHeader.CurHash), block.Data, blockchain.NewProofOfWork(block, blockchain.POW_TARGET).Validate())
	}
}


func (cli *Cli) PrintBlock(block *blockchain.Block) {
	fmt.Printf(" %x , %x,  %s, %b \n",string(block.BlockHeader.PrevHash), string(block.BlockHeader.CurHash), block.Data, blockchain.NewProofOfWork(block, blockchain.POW_TARGET).Validate())
}


func (cli *Cli) ValidateBlock(b * blockchain.Block) bool {
	return blockchain.NewProofOfWork(b, blockchain.POW_TARGET).Validate()
}

