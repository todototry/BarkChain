package blockchain

import (
	"time"
	"strconv"
	"bytes"
	"crypto/sha256"
)


type BlockHeader struct {
	TimeStamp int64
	PrevHash []byte
	Nonce  int
	CurHash []byte
}


type Block struct {
	BlockHeader *BlockHeader
	Data string
}


type BlockChain struct {
	tip *Block
	Blocks []*Block
}


func (block *Block) SetHashAuto() {

	timestamp := []byte(strconv.FormatInt(block.BlockHeader.TimeStamp, 10))
	headers := bytes.Join([][]byte{block.BlockHeader.PrevHash, []byte(block.Data), timestamp}, []byte{})
	hash := sha256.Sum256(headers)

	block.BlockHeader.CurHash = hash[:]
}


func (block *Block) PrepareData() []byte {
	timestamp := []byte(strconv.FormatInt(block.BlockHeader.TimeStamp, 10))
	nonce := []byte(strconv.Itoa(block.BlockHeader.Nonce))
	data := []byte(block.Data)

	food := bytes.Join([][]byte{timestamp, block.BlockHeader.PrevHash, nonce, data}, []byte{})

	return food
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
	// block.SetHashAuto()

	proof := NewProofOfWork(block, POW_TARGET)
	proof.Pow()

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
	// b.SetHashAuto()

	proof := NewProofOfWork(b, POW_TARGET)
	proof.Pow()

	bc.Blocks = append(bc.Blocks, b)
}
