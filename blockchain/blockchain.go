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


type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

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

	// create bolt Db
	db, err := bolt.Open(DB_PATH, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err == nil {
		bc.Db = db
	} else {
		os.Exit(-1)
	}

	// if bucket Not exists save genesis blockGenesis to DB
	bc.Db.Update(func(tx *bolt.Tx) error {
		bucket:= tx.Bucket([]byte(BUCKET_NAME))

		// if block data does not exist in DB
		if bucket == nil {
			bucket,err = tx.CreateBucket([]byte(BUCKET_NAME))

			if err == nil {
				// create and append GenesisBlock to bolt DB
				blockGenesis := newGenesisBlock()
				bc.Blocks = append(bc.Blocks, blockGenesis)
				// update Tip
				bc.Tip = blockGenesis.BlockHeader.CurHash

				// save as kv after block serialized.
				err = bucket.Put(blockGenesis.BlockHeader.CurHash, blockGenesis.Serialize())

				// save tip
				err = bucket.Put([]byte(BLOCKCHAIN_TIP_KEY), blockGenesis.BlockHeader.CurHash)
			}
		} else {
			// blocks already exists in DB.
			bc.Tip = bucket.Get([]byte(BLOCKCHAIN_TIP_KEY))
		}

		return err
	})

	return bc
}


// new block for the blockchain.
func (bc *BlockChain) NewBlock(data string)  {

	b := &Block{
		BlockHeader: &BlockHeader{
			CurHash:nil,
			PrevHash:bc.Tip,
			TimeStamp: time.Now().Unix(),
		},
		Data:data,
	}
	// cal hash for pow
	proof := NewProofOfWork(b, POW_TARGET)
	proof.Pow()

	// update blockChain.Tip in mem.
	bc.Tip = b.BlockHeader.CurHash

	// save to bolt DB.
	bc.Db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(BUCKET_NAME))

		// 1. save as kv after block serialized  to bolt DB.
		err := tx.Bucket([]byte(BUCKET_NAME)).Put(b.BlockHeader.CurHash, b.Serialize())

		// 2. update Tip in bolt DB.
		tx.Bucket([]byte(BUCKET_NAME)).Put([]byte(BLOCKCHAIN_TIP_KEY), b.BlockHeader.CurHash)

		return err
	})

}


func (bc *BlockChain) Iterator() *BlockChainIterator {
	bci:= &BlockChainIterator{bc.Tip, bc.Db}
	return bci
}


type BlockChainIterator struct {
	CurHash []byte
	Db *bolt.DB
}


// iterate the blockchain.
func (bci *BlockChainIterator) Next() *Block {
	var b *Block

	err := bci.Db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(BUCKET_NAME))

		if bucket == nil{
			// fail
			return &errorString{"Bucket not exists..."}
		} else {
			k := bci.CurHash
			// key is nil, indicates that finished iterate.
			if k == nil {
				b = nil
				return nil
			} else {
				// get data and deserialize it.
				v := bucket.Get(k)
				if v == nil {
					b = nil
				} else {
					b = Deserialize(v)
				}
			}

			return nil
		}
	})

	if err == nil {
		// success
		if b == nil {
			return nil
		} else {
			bci.CurHash = b.BlockHeader.PrevHash
			return 	b
		}
	} else {
		return nil
	}
}

