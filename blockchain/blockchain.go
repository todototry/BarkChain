package blockchain

import (
	"time"
	"strconv"
	"bytes"
	"crypto/sha256"
	"github.com/boltdb/bolt"
	"os"
	"encoding/gob"
	"log"
)

const (
	DB_PATH = "btc_blockchain.Db"
	BUCKET_NAME= "BTC_BLOCKCHAIN"
	BLOCKCHAIN_TIP_KEY = "TIP_OF_BLOCKCHAIN"
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
	Tip    []byte
	Blocks []*Block
	Db     *bolt.DB
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


func (block *Block) Serialize() []byte {
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)

	err := encoder.Encode(block)

	if err != nil {
		log.Panic(err)
	}

	return buf.Bytes()
}


func Deserialize(buf []byte) *Block {
	var b Block
	var data = bytes.NewReader(buf)

	decoder := gob.NewDecoder(data)

	err := decoder.Decode(&b)
	if err != nil {
		log.Panic(err)
	}

	return &b
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

	// save to Bolt DB.
	return block
}


func NewBlockChain() *BlockChain {
	var bc = &BlockChain{}

	// append GenesisBlock
	blockGenesis := newGenesisBlock()
	bc.Blocks = append(bc.Blocks, blockGenesis)
	// update Tip
	bc.Tip = blockGenesis.BlockHeader.CurHash


	// create bolt Db
	db, err := bolt.Open(DB_PATH, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err == nil {
		bc.Db = db
	} else {
		os.Exit(-1)
	}

	// save genesis blockGenesis to DB
	bc.Db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BUCKET_NAME))

		// save as kv after block serialized.
		err = bucket.Put(blockGenesis.BlockHeader.CurHash, blockGenesis.Serialize())
		return err
	})

	// save Tip
	bc.Db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(BUCKET_NAME))

		// save tip
		err = bucket.Put([]byte(BLOCKCHAIN_TIP_KEY), blockGenesis.BlockHeader.CurHash)
		return err
	})

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
	bc.Tip = b.BlockHeader.CurHash

	// save to bolt DB.
	bc.Db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(BUCKET_NAME))

		// save as kv after block serialized.
		err := tx.Bucket([]byte(BUCKET_NAME)).Put(b.BlockHeader.CurHash, b.Serialize())

		// update Tip
		tx.Bucket([]byte(BUCKET_NAME)).Put([]byte(BLOCKCHAIN_TIP_KEY), b.BlockHeader.CurHash)

		return err
	})

	bc.Db.Sync()
}
