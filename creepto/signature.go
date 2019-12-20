package creepto

import (
	"encoding/asn1"
	"encoding/hex"
	"math/big"
)

//Signature represents a string
type Signature struct {
	R *big.Int
	S *big.Int
}

//Generates the DER of the signature
func (sig *Signature) DER() string {
	encoded, _ := asn1.Marshal(*sig)
	return hex.EncodeToString(encoded)
}

//Gets the Signature from the DER
