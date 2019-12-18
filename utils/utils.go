package utils

import (
	"crypto/rand"
	"math/big"
)

//Creates a determinitic number
func GenerateSecret() string {
	//Max random value, a 130-bits integer, i.e 2^130 - 1
	max := new(big.Int)
	max.Exp(big.NewInt(2), big.NewInt(130), nil).Sub(max, big.NewInt(1))

	//Generate cryptographically strong pseudo-random between 0 - max
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		//error handling
	}
	//String representation of n in base 16
	return n.Text(16)

}
