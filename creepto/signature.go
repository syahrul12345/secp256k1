package creepto

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/syahrul12345/secp256k1/utils"
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
	rbin, _ := hex.DecodeString(sig.R.Text(16))
	for rbin[0] == 0 {
		// Strip leading 0 bytes
		rbin = rbin[1:]
	}
	if rbin[0]&0x80 == 1 {
		// Append a 0 byte
		rbin = append([]byte{0}, rbin...)
	}
	result := append([]byte{
		2,
		byte(len(rbin)),
	}, rbin...)
	// repeat the same for sbin
	sbin, _ := hex.DecodeString(sig.S.Text(16))
	for sbin[0] == 0 {
		//Strip leading 0 bytes
		sbin = sbin[1:]
	}
	if sbin[0]&0x80 == 1 {
		// Append a 0 byte
		sbin = append([]byte{0}, rbin...)
	}
	sbinResult := append([]byte{
		2,
		byte(len(sbin)),
	}, sbin...)
	result = append(result, sbinResult...)
	final := append([]byte{
		0x30,
		byte(len(result)),
	}, result...)
	return hex.EncodeToString(final)
	// encoded, _ := asn1.Marshal(*sig)
	// return hex.EncodeToString(encoded)
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
func ParseDER(derstring string) *Signature {
	buf, err := hex.DecodeString(derstring)
	bufLength := len(buf)
	if err != nil {
		fmt.Println("Failed to serialize DER string")
	}
	// Read the first object
	compound := buf[0]
	buf = buf[1:]
	if compound != 0x30 {
		fmt.Println("Bad Signature")
		return nil
	}
	length := int(buf[0])
	buf = buf[1:]
	if length+2 != bufLength {
		fmt.Println("Bad signature length")
		return nil
	}
	marker := buf[0]
	buf = buf[1:]
	if marker != 0x02 {
		fmt.Println("Bad signaure")
		return nil
	}

	rlength := int8(buf[0])
	buf = buf[1:]
	r := big.NewInt(0).SetBytes(buf[:rlength])
	buf = buf[rlength:]

	marker = buf[0]
	buf = buf[1:]
	if marker != 0x02 {
		fmt.Println("bad Signature")
		return nil
	}
	slength := int8(buf[0])
	buf = buf[1:]
	s := big.NewInt(0).SetBytes(buf[:slength])
	if len(derstring)/2 != int(6+rlength+slength) {
		fmt.Println("signature is too long")
		return nil
	}
	return &Signature{
		r,
		s,
	}
}
