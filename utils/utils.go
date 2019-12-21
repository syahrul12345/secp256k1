package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"golang.org/x/crypto/ripemd160"
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

//Encode58 will convert hexadecimal strings to base58
func Encode58(str string) string {
	hexBytes, _ := hex.DecodeString(str)
	encoded := base58.Encode(hexBytes)
	return encoded

}

//Encode58CheckSum will include the checksum to a base58 encoded string
func Encode58CheckSum(str string) string {
	hexBytes, _ := hex.DecodeString(str)
	preCheckSum := sha256.Sum256(hexBytes)
	checksum := sha256.Sum256(preCheckSum[:])
	first4HexString := hex.EncodeToString(checksum[0:4])
	return Encode58(str + first4HexString)
}

//Hash160 utility will take the SEC address format, applies a sha256 then ripemd160
func Hash160(sec string) string {
	hexBytes, _ := hex.DecodeString(sec)
	hash256 := sha256.Sum256(hexBytes)
	ripemdHasher := ripemd160.New()
	ripemdHasher.Write(hash256[:])
	hashBytes := ripemdHasher.Sum(nil)
	hashString := fmt.Sprintf("%x", hashBytes)
	return hashString
}

func divmod(numerator, denominator uint64) (quotient, remainder uint64) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}
