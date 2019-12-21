package creepto

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"secp256k1/curve"
	"secp256k1/utils"
)

//PrivateKey represents a private key that can be used to decrypt messages.
//Stores an exportable PublicKey as a *Point256
//Stores an unexportable secret. This prevents unwanted dumping
type PrivateKey struct {
	PublicKey *Point256
	secret    string
}

//CreateNewPrivateKey will create a new private key object that can be used to sign
//Can be triggered from the console
func CreateNewPrivateKey() PrivateKey {
	//Generate a random secret
	G, _ := curve.NewPoint("0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", "0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8")
	//In production we generate a deterministic number
	//Create a hex of the secret
	secret := utils.GenerateSecret()
	pubKey, _ := G.Mul("0x" + secret)
	pubKey256 := New256Point(utils.FromBigInt(pubKey.X.Number), utils.FromBigInt(pubKey.Y.Number))
	return PrivateKey{
		pubKey256,
		secret,
	}
}

//DumpSecret will display the PrivateKey in a human readable format
func (privKey PrivateKey) DumpSecret() string {
	return privKey.secret
}

//GetPublicKey will derive a Public Key given a private key.
func GetPublicKey(secret string) *Point256 {
	//Generator Point for bitcoin
	G, _ := curve.NewPoint("0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", "0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8")
	pubKey, _ := G.Mul(secret)
	return New256Point(utils.FromBigInt(pubKey.X.Number), utils.FromBigInt(pubKey.Y.Number))
}

//Sign willl sign a message with a privateKey
func (privKey PrivateKey) Sign(message string) (*Signature, string) {
	//Generator point
	G, _ := curve.NewPoint("0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", "0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8")
	//Convert and hashes the string to be signed
	var data []byte = []byte(message)
	hasher := sha256.New()
	hasher.Write(data)
	z := hex.EncodeToString(hasher.Sum(nil))
	//End of hashing phase

	//Calculate the required mathematical constants. Generate a deterministic k
	k := utils.GenerateSecret()
	R, _ := G.Mul("0x" + k)
	bigOrder := utils.ToBigInt(Order)
	newOrder := big.NewInt(0).Sub(bigOrder, big.NewInt(2))
	bigK := utils.ToBigInt("0x" + k)
	bigZ := utils.ToBigInt("0x" + z)
	bigSecret := utils.ToBigInt("0x" + privKey.secret)
	//Do the actual math to generate signatures
	//Calculate: 1/k
	kInv := big.NewInt(0).Exp(bigK, newOrder, bigOrder)
	//Calculate r*secret
	term1 := big.NewInt(0).Mul(R.X.Number, bigSecret)
	//Calculate z + term1
	term2 := big.NewInt(0).Add(bigZ, term1)
	//Calculate term2 * kInv
	term3 := big.NewInt(0).Mul(term2, kInv)
	//Calculate term3 %n
	s := big.NewInt(0).Mod(term3, bigOrder)
	if s.Cmp(big.NewInt(0).Div(bigOrder, big.NewInt(2))) == 1 {
		s = big.NewInt(0).Sub(bigOrder, s)
	}
	return &Signature{
		R.X.Number,
		s,
	}, "0x" + z
}

//Export will display the privatekey as WIF format
func (privKey PrivateKey) Export(compressed bool, testnet bool) string {
	secret := privKey.DumpSecret()
	var prefix string
	var suffix string
	if testnet {
		prefix = "6f"
	} else {
		prefix = "00"
	}
	if compressed {
		suffix = "01"
	} else {
		suffix = ""
	}
	return utils.Encode58CheckSum(prefix + secret + suffix)
}
