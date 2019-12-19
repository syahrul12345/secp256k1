package creepto

import (
	"fmt"
	"math/big"
	"secp256k1/curve"
	"secp256k1/fieldelement"
	"secp256k1/utils"
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

//Converts a normal Point to a Point256
func NormalTo256(point *curve.Point) *Point256 {
	return &Point256{
		point.X,
		point.Y,
		point.A,
		point.B,
	}
}

//Mul will multiply point256 with a coefficient
func (point256 *Point256) Mul(coefficient string) (*curve.Point, error) {
	// Check if the coefifcinent is already hexed
	coeff := utils.ToBigInt(coefficient)
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

//SEC will create the serialized public key for propogation
//Takes a boolean paramenter to determine if the output is compressed or not
func (point256 *Point256) SEC(compressed bool) string {
	if compressed {
		x := big.NewInt(0).Mod(point256.X.Number, big.NewInt(2))
		xBytes := new(big.Int).SetBytes(point256.X.Number.Bytes())
		text := xBytes.Text(16)
		if x.Cmp(big.NewInt(0)) == 0 {
			return "02" + text
		}
		return "03" + text
	}
	xBytes := new(big.Int).SetBytes(point256.X.Number.Bytes())
	yBytes := new(big.Int).SetBytes(point256.Y.Number.Bytes())
	textX := xBytes.Text(16)
	textY := yBytes.Text(16)
	return "04" + textX + textY
}

//ParseBin will discover the public key given a SEC_BIN
func ParseBin(secBin string) *Point256 {
	// Check if it's compressed or not
	flag, _ := strconv.ParseInt(secBin[1:2], 10, 64)
	if flag == 4 {
		xHalf := secBin[2:66]
		yHalf := secBin[66:]
		return New256Point("0x"+xHalf, "0x"+yHalf)
	}
	xHalf := secBin[2:]
	xHex := "0x" + xHalf
	// Create the required fields
	xField := fieldelement.NewFieldElement(xHex)
	alpha := xField.Pow("3").Add(fieldelement.NewFieldElement("7"))
	beta := alpha.Sqrt()
	//Check if 0 == beta % 2 and check if it's 0 a
	//Double check 0 as the big library returns a 0
	var evenBeta fieldelement.FieldElement
	var oddBeta fieldelement.FieldElement
	var input *big.Int
	input = big.NewInt(0).Sub(beta.Prime, beta.Number)
	if big.NewInt(0).Cmp(big.NewInt(0).Mod(beta.Number, big.NewInt(2))) == 0 {
		evenBeta = beta
		oddBeta = fieldelement.NewFieldElement(input.Text(10))
	} else {
		evenBeta = fieldelement.NewFieldElement(input.Text(10))
		oddBeta = beta

	}
	//even
	if flag == 2 {
		return New256Point(xHex, evenBeta.Number.Text(10))
	} else {
		return New256Point(xHex, oddBeta.Number.Text(10))
	}

}
