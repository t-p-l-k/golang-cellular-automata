package main

import (
	"crypto/rand"
	"math/big"
	"fmt"
)

func GenerateRandomUint64(max int64) uint64 {
	n, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		fmt.Println(err)
	}
	return n.Uint64()
}
