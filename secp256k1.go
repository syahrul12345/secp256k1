package secp256k1

import (
	"github.com/syahrul12345/secp256k1/creepto"
	"github.com/syahrul12345/secp256k1/utils"
)

//NewBitCoinAddress  will generate a bitcoin address
func NewBitCoinAddress(compressed bool, testnet bool) (pubKey string, secret string) {
	privateKey := creepto.CreateNewPrivateKey()
	pubKey = privateKey.PublicKey.GetAddress(compressed, testnet)
	secret = privateKey.DumpSecret()
	return
}

//GetMainnetAddressFromPrivKey : Returns the public key for a given private key in the mainnet
func GetMainnetAddressFromPrivKey(secret string) string {
	point256 := creepto.GetPublicKey(secret)
	return point256.GetAddress(true, false)
}

//GetSec: Gets the SEC format of a public key, given a private key
func GetSec(secret string, compressed bool) string {
	point256 := creepto.GetPublicKey(secret)
	return point256.SEC(compressed)
}

//DecodeAddress: Decodes an address, and return the hashed sec string.
func DecodeAddress(address string) string {
	return utils.Decode58(address)
}

//GetTestnetAddressFromPrivKey : Returns the public key for a given private key in the testnet
func GetTestnetAddressFromPrivKey(secret string) string {
	point256 := creepto.GetPublicKey(secret)
	return point256.GetAddress(true, true)
}

//Signs a message given a privatekey
func Sign(secret string, message string) (*creepto.Signature, string) {
	privKeyObj := creepto.CreateNewPrivateKeyFromSecret(secret)
	return privKeyObj.Sign(message)
}

//Parse the sec in the string format and return the corresponding point256
func ParseSec(secPubKey string) *creepto.Point256 {
	return creepto.ParseSec(secPubKey)
}

//Parse the derString and returns the unserialized Signature Object
func ParseDer(derString string) *creepto.Signature {
	return creepto.ParseDER(derString)
}

//Verify : Verifies if the der signature string, a signature hash is sent by the sec PUbkey
func Verify(secPubKey string, der string, z string) (bool, error) {
	point256 := creepto.ParseSec(secPubKey)
	signature := creepto.ParseDER(der)
	return point256.Verify(z, signature)
}
