package utils

import (
	"crypto/rand"
	"math/big"
)

//GenerateSecret will create  a deterministic number
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

//ToBigInt Converts a string to the corresponding big.Int
func ToBigInt(ridiculouslyLargeNumber string) *big.Int {
	var decodedNumber *big.Int
	// First handle nil values
	if ridiculouslyLargeNumber == "nil" {
		return nil
	}
	// Check for length, if more than 3, it's definately a hex
	if len(ridiculouslyLargeNumber) >= 3 {
		hexCheck := ridiculouslyLargeNumber[0:2]
		isHex := (hexCheck == "0x")
		if isHex {
			//Remove the 0x prefix
			decodedNumber, _ = big.NewInt(0).SetString(ridiculouslyLargeNumber[2:], 16)
		} else {
			decodedNumber, _ = big.NewInt(0).SetString(ridiculouslyLargeNumber, 10)
		}
	} else {
		// Not a hex, so we can set the string
		decodedNumber, _ = big.NewInt(0).SetString(ridiculouslyLargeNumber, 10)

	}

	return decodedNumber
}

//FromBigInt will return the 0x encoded of the provided bigInt
func FromBigInt(bigInt *big.Int) string {
	return "0x" + bigInt.Text(16)
}
