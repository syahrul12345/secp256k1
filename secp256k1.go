package secp256k1

import (
	"secp256k1/creepto"
)

//NewBitCoinAddress  will generate a bitcoin address
func NewBitCoinAddress(compressed bool, testnet bool) (pubKey string, secret string) {
	privateKey := creepto.CreateNewPrivateKey()
	pubKey = privateKey.PublicKey.GetAddress(compressed, testnet)
	secret = privateKey.DumpSecret()
	return
}
