package creepto

import (
	"fmt"
	"math/big"
	"secp256k1/curve"
	"secp256k1/fieldelement"
	"strconv"

	"github.com/bitherhq/go-bither/common/hexutil"
)

const (
	//Order of fin
	Order string = "0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141"
)

//Point256 represents the public key on the elliptical curve
type Point256 struct {
	X *fieldelement.FieldElement
	Y *fieldelement.FieldElement
	A fieldelement.FieldElement
	B fieldelement.FieldElement
}

//Signature represents a string
type Signature struct {
	R string
	S string
}

//New256Point creates a Point256 representing a public key
func New256Point(x string, y string) *Point256 {
	newPoint, pointError := curve.NewPoint(x, y)
	if pointError != nil {
		fmt.Println(pointError)
		return nil
	}
	return &Point256{
		newPoint.X,
		newPoint.Y,
		newPoint.A,
		newPoint.B,
	}
}

//Mul will multiply point256 with a coefficient
func (point256 *Point256) Mul(coefficient string) (*curve.Point, error) {
	// Check if the coefifcinent is already hexed
	_, alreadyHex := strconv.ParseInt(coefficient, 16, 64)
	var coeff *big.Int
	if alreadyHex != nil {
		tempCoeff, _ := hexutil.DecodeBig(coefficient)
		coeff = tempCoeff
	} else {
		coefficientInterger, _ := strconv.Atoi(coefficient)
		hexString := fmt.Sprintf("%x", coefficientInterger)
		tempCoeff, _ := hexutil.DecodeBig("0x" + hexString)
		coeff = tempCoeff
	}

	modoBig, err := hexutil.DecodeBig(Order)
	if err != nil {
		return nil, nil
	}
	//Sets the new coefficient
	coeff.Mod(coeff, modoBig)
	coeffString := hexutil.EncodeBig(coeff)
	// We create a normal point that can do the multiplacaiton
	x := point256.X.Number
	y := point256.Y.Number
	//Skip error handling, it will definately be on the line
	tempPoint, _ := curve.NewPoint(hexutil.EncodeBig(x), hexutil.EncodeBig(y))
	result, _ := tempPoint.Mul(coeffString)

	//Return and reconvert it to  a point256
	return result, nil
}

//NewSignature will create a new Signature
func NewSignature(r string, s string) *Signature {
	return &Signature{
		r,
		s,
	}
}

//Verify This function will verify if the Public Key has sent the message z, with enclosed signature
func (point256 *Point256) Verify(message string, sig *Signature) bool {
	z, _ := hexutil.DecodeBig(message)
	r, _ := hexutil.DecodeBig(sig.R)
	s, _ := hexutil.DecodeBig(sig.S)

	bigOrder, _ := hexutil.DecodeBig(Order)
	newOrder := big.NewInt(0).Sub(bigOrder, big.NewInt(2))
	sInv := big.NewInt(0).Exp(s, newOrder, bigOrder)
	//u = z/s
	u1 := big.NewInt(0).Mul(z, sInv)
	u := big.NewInt(0).Mod(u1, bigOrder)
	//v = r/s
	v1 := big.NewInt(0).Mul(r, sInv)
	v := big.NewInt(0).Mod(v1, bigOrder)
	//R(r,s) = uG + vP
	//Create G
	G, _ := curve.NewPoint("0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", "0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8")

	firstTerm, _ := G.Mul(hexutil.EncodeBig(u))
	secondTerm, _ := point256.Mul(hexutil.EncodeBig(v))
	R, err := firstTerm.Add(secondTerm)
	if err != nil {
		fmt.Println(err)
	}

	res := R.X.Number.Cmp(r)
	if res == 0 {
		return true
	}
	return false
}
