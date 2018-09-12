package blockchain

import (
	"math/big"
	"math"
	"crypto/sha256"
	"fmt"
)

const (

	POW_TARGET = 16
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
		fmt.Printf("\r %v ", hash)

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

func (proof *ProofOfWork) Validate() bool {
	// prepare
	data := proof.block.PrepareData()

	// calc sha
	hash := sha256.Sum256(data)

	// check
	resInt := big.NewInt(1).SetBytes(hash[:])
	if resInt.Cmp(proof.target) < 0 {
		return true
	} else {
		return false
	}
}
