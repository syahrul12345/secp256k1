package creepto

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"secp256k1/curve"
	"secp256k1/utils"

	"github.com/bitherhq/go-bither/common/hexutil"
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
	secret := utils.GenerateSecret()
	secretBig, _ := hexutil.DecodeBig("0x" + secret)
	pubKey, _ := G.Mul(hexutil.EncodeBig(secretBig))
	pubKey256 := New256Point(hexutil.EncodeBig(pubKey.X.Number), hexutil.EncodeBig(pubKey.Y.Number))
	return PrivateKey{
		pubKey256,
		secret,
	}
}

//GetPublicKey will derive a Public Key given a private key.
func GetPublicKey(secret string) *Point256 {
	//Generator Point for bitcoin
	G, _ := curve.NewPoint("0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", "0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8")
	pubKey, _ := G.Mul(secret)
	return New256Point(hexutil.EncodeBig(pubKey.X.Number), hexutil.EncodeBig(pubKey.Y.Number))
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
	bigOrder, _ := hexutil.DecodeBig(Order)
	newOrder := big.NewInt(0).Sub(bigOrder, big.NewInt(2))
	bigK, _ := hexutil.DecodeBig("0x" + k)
	bigZ, _ := hexutil.DecodeBig("0x" + z)
	bigSecret, _ := hexutil.DecodeBig("0x" + privKey.secret)
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
		hexutil.EncodeBig(R.X.Number),
		hexutil.EncodeBig(s),
	}, "0x" + z

}
