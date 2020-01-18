package creepto

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"math/big"

	"github.com/syahrul12345/secp256k1/curve"
	"github.com/syahrul12345/secp256k1/utils"
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

//CreateNewPrivateKeyFromSecret will create a new privkey from a secret
func CreateNewPrivateKeyFromSecret(secret string) PrivateKey {
	pubKey := GetPublicKey(secret)
	return PrivateKey{
		pubKey,
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

//Sign willl sign a message with a privateKey. This will return a signature object, and the signature hash
func (privKey PrivateKey) Sign(message string) (*Signature, string) {
	//Generator point
	G, _ := curve.NewPoint("0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", "0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8")
	//Convert and hashes the string to be signed
	var data []byte = []byte(message)
	hasher := sha256.New()
	hasher.Write(data)
	sigHash := hex.EncodeToString(hasher.Sum(nil))
	//End of hashing phase
	//Calculate the required mathematical constants. Generate a deterministic k
	k := privKey.GenerateDeterministicK(message)
	R, _ := G.Mul("0x" + k)
	bigOrder := utils.ToBigInt(Order)
	newOrder := big.NewInt(0).Sub(bigOrder, big.NewInt(2))
	bigK := utils.ToBigInt("0x" + k)
	bigZ := utils.ToBigInt("0x" + message)
	bigSecret := utils.ToBigInt(privKey.DumpSecret())
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
	}, "0x" + sigHash
}

func digest(key []byte, buf []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(buf)
	return mac.Sum(nil)
}

// Generate a determistic number,based on the input
func (privkey PrivateKey) GenerateDeterministicK(input string) string {
	Order := "0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141"
	OrderBig, _ := big.NewInt(0).SetString(Order[2:], 16)
	z, _ := big.NewInt(0).SetString(input, 16)
	if z.Cmp(OrderBig) == 1 {
		z.Sub(z, OrderBig)
	}
	k := []byte{}
	v := []byte{}
	// Fill up the bytes
	counter := 0
	for counter < 32 {
		k = append(k, byte(0))
		v = append(v, byte(01))
		counter++
	}
	zBytes, _ := hex.DecodeString(input)
	// privkey secret is in a string... not even a hexadecimal string
	secretBytes := utils.ToBigInt(privkey.DumpSecret()).Bytes()
	// We want to pad it with 0 bytes until it reaches 32 bytes
	for len(secretBytes) < 32 {
		secretBytes = append([]byte{0x00}, secretBytes...)
	}
	var buf []byte

	routes := [][]byte{
		v,
		[]byte{0x00},
		secretBytes,
		zBytes,
	}
	for _, route := range routes {
		buf = append(buf, route...)
	}
	k = digest(k, buf)
	v = digest(k, v)
	buf = []byte{}
	routes = [][]byte{
		v,
		[]byte{0x01},
		secretBytes,
		zBytes,
	}
	for _, route := range routes {
		buf = append(buf, route...)
	}
	k = digest(k, buf)
	v = digest(k, v)
	for true {
		v = digest(k, v)
		candidate := big.NewInt(0).SetBytes(v)
		if candidate.Cmp(big.NewInt(1)) == 1 && candidate.Cmp(OrderBig) == -1 {
			return candidate.Text(16)
		}
		buf = []byte{}
		routes = [][]byte{
			v,
			[]byte{0x00},
		}
		for _, route := range routes {
			buf = append(buf, route...)
		}
		k = digest(k, buf)
		v = digest(k, v)
	}
	return "s"
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

//GetAddress will get teh address of the privatekey
func (privKey PrivateKey) GetAddress(compressed bool, testnet bool) string {
	return privKey.PublicKey.GetAddress(compressed, testnet)
}
