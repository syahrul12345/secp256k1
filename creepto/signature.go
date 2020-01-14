package creepto

import (
	"encoding/asn1"
	"encoding/hex"
	"github.com/syahrul12345/secp256k1/utils"
	"math/big"
)

//Signature represents a string
type Signature struct {
	R *big.Int
	S *big.Int
}

//NewSignature creates an arbitary new signature
func NewSignature(r, s string) *Signature {
	return &Signature{
		utils.ToBigInt(r),
		utils.ToBigInt(s),
	}
}

//Generates the DER of the signature
func (sig *Signature) DER() string {
	encoded, _ := asn1.Marshal(*sig)
	return hex.EncodeToString(encoded)
}

//Equals check if two signatures are the same
func (sig *Signature) Equals(sig2 *Signature) bool {
	rEquals := sig.R.Cmp(sig2.R)
	sEquals := sig.S.Cmp(sig.S)
	if rEquals == 0 && sEquals == 0 {
		return true
	}
	return false
}

//ParseDER will accept a DER serialized signature, and obtaine the unserialized struct
func ParseDER(s string) *Signature {
	tempSig := new(Signature)
	derByte, _ := hex.DecodeString(s)
	asn1.Unmarshal(derByte, tempSig)
	return tempSig
}
