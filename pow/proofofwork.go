package main

import (
	"math/big"
	"math"
	"crypto/sha256"
	"bytes"
	"strconv"
	"fmt"
)

func pow(data []byte, targetbits int) (int, [32]byte){
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
	pow([]byte(" I am fine, Thank you!"), 64)
}
