package main

import (
	"fmt"
	"secp256k1/creepto"
)

//NewBitCoinAddress  will generate a bitcoin address
func NewBitCoinAddress(compressed bool, testnet bool) (pubKey string, secret string) {
	privateKey := creepto.CreateNewPrivateKey()
	pubKey = privateKey.PublicKey.GetAddress(compressed, testnet)
	secret = privateKey.DumpSecret()
	return
}

func main() {
	pubKey := creepto.GetPublicKey("0xdeadbeef")
	secCorrect := pubKey.SEC(true)
	res := creepto.ParseSec(secCorrect)
	fmt.Println(res.SEC(true))

}
