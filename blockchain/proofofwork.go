package blockchain

import (
	"math/big"
	"math"
	"crypto/sha256"
	"bytes"
	"strconv"
	"fmt"
)

type ProofOfWork struct {
	block *Block
	target *big.Int
}

func NewProofOfWork(block *Block, targetbits int) *ProofOfWork {

	target := big.NewInt(1)
	target.Lsh(target, uint(256 - targetbits))

	pow := ProofOfWork{block, target}

	return &pow
}


func (proof * ProofOfWork) Pow() (nonce int, hash [32]byte){

	for nonce = 0; nonce <= math.MaxInt32; nonce ++ {
		// set nonce
		proof.block.BlockHeader.Nonce = nonce

		// get data
		data := proof.block.PrepareData()

		// calc sum256
		hash := sha256.Sum256(data)
		fmt.Printf("\r %v", hash)

		// check valid, compare with target.
		resInt := big.NewInt(1).SetBytes(hash[:])

		if resInt.Cmp(proof.target) < 0 {
			// found and set.
			proof.block.BlockHeader.CurHash = hash[:]
			break
		}
	}

	return nonce, hash
}


func Pow(data []byte, targetbits int) (int, [32]byte){
	max_nonce := math.MaxInt64
	hash_result := big.NewInt(1)

	// 1. target value
	origin := big.NewInt(1)
	targetvalue := origin.Lsh(origin, uint(256-targetbits))

	// 2. find a nonce, who's sum256 is less than target value
	nonce :=0
	var dig [32]byte

	for ; nonce <= max_nonce; nonce++ {
		// 2.1 random nonce

		// 2.2 prepare data + nonce
		food := bytes.Join(
			[][]byte{
				data,
				[]byte(strconv.Itoa(nonce)),
			},
			[]byte{},
			)

		// 2.3 calc sum256
		dig = sha256.Sum256(food)
		fmt.Printf("\r %x", dig)
		// 2.4 compare sum256 with target value.
		if targetvalue.Cmp(hash_result.SetBytes(dig[:])) > 0 {
			break
		}
	}

	// 3. return nonce, sum256
	return nonce, dig
}


func main() {
	Pow([]byte(" I am fine, Thank you!"), 24)
}
