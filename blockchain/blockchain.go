package blockchain

import (
	"time"
	"strconv"
	"bytes"
	"crypto/sha256"
)

type BlockHeader struct {
	CurHash []byte
	PrevHash []byte
	TimeStamp int64
}


type Block struct {
	BlockHeader *BlockHeader
	Data string
}


type BlockChain struct {
	tip *Block
	Blocks []*Block
}


func (b *Block) SetHash() {

	timestamp := []byte(strconv.FormatInt(b.BlockHeader.TimeStamp, 10))
	headers := bytes.Join([][]byte{b.BlockHeader.PrevHash, []byte(b.Data), timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	b.BlockHeader.CurHash = hash[:]
}


func newGenesisBlock() *Block {

	var block = &Block{
		BlockHeader:&BlockHeader{
			CurHash:nil,
			PrevHash:nil,
			TimeStamp:time.Now().Unix(),
		},
		Data:"Genesis Block",
	}

	block.SetHash()

	return block
}


func NewBlockChain() *BlockChain {
	var bc = &BlockChain{}

	// append GenesisBlock
	block := newGenesisBlock()
	bc.Blocks = append(bc.Blocks, block)
	return bc
}


// new block for the blockchain.
func (bc *BlockChain) NewBlock(data string)  {

	lastIndex := len(bc.Blocks) - 1
	lastBlock := bc.Blocks[lastIndex]

	b := &Block{
		BlockHeader: &BlockHeader{
			CurHash:nil,
			PrevHash:nil,
			TimeStamp: time.Now().Unix(),
		},
		Data:data,
	}

	b.BlockHeader.PrevHash = lastBlock.BlockHeader.CurHash
	b.SetHash()

	bc.Blocks = append(bc.Blocks, b)
}
