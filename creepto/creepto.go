package creepto

import (
	"encoding/hex"
	"fmt"
	"github.com/syahrul12345/secp256k1/curve"
	"github.com/syahrul12345/secp256k1/fieldelement"
	"github.com/syahrul12345/secp256k1/utils"
	"math/big"
)

var (
	//Order of fin
	Order string = "0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141"
	//P is the max prime value
	P string = big.NewInt(0).Sub(
		big.NewInt(0).Sub(
			big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0)),
			big.NewInt(0).Exp(big.NewInt(2), big.NewInt(32), big.NewInt(0))),
		big.NewInt(977)).String()
)

//Point256 represents the public key on the elliptical curve
type Point256 struct {
	X *fieldelement.FieldElement
	Y *fieldelement.FieldElement
	A fieldelement.FieldElement
	B fieldelement.FieldElement
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

//NormalTo256 will convert the given point to a Point256 type
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
	modoBig, ok := big.NewInt(0).SetString(Order[2:], 16)
	if !ok {
		return nil, nil
	}
	//Sets the new coefficient

	coeff.Mod(coeff, modoBig)
	coeffString := "0x" + coeff.Text(16)
	// We create a normal point that can do the multiplacaiton
	x := point256.X.Number
	y := point256.Y.Number
	//Skip error handling, it will definately be on the line
	tempPoint, _ := curve.NewPoint("0x"+x.Text(16), "0x"+y.Text(16))
	result, _ := tempPoint.Mul(coeffString)

	//Return and reconvert it to  a point256
	return result, nil
}

//Verify This function will verify if the Public Key has sent the signature hash z, with enclosed signature
func (point256 *Point256) Verify(signatureHash string, sig *Signature) bool {
	z := utils.ToBigInt(signatureHash)
	r := sig.R
	s := sig.S

	bigOrder := utils.ToBigInt(Order)
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

	firstTerm, _ := G.Mul("0x" + u.Text(16))
	secondTerm, _ := point256.Mul("0x" + v.Text(16))
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

//ParseSec will return the point256
func ParseSec(secString string) *Point256 {
	// Lets convert the string to a bytes
	byteSec, _ := hex.DecodeString(secString)
	if byteSec[0] == 4 {
		x := "0x" + hex.EncodeToString(byteSec[1:33])
		y := "0x" + hex.EncodeToString(byteSec[33:65])
		return New256Point(x, y)
	}
	var isEven bool
	if byteSec[0] == 2 {
		isEven = true
	} else {
		isEven = false
	}
	bytesString := "0x" + hex.EncodeToString(byteSec[1:])
	x := fieldelement.NewFieldElement(bytesString)
	a := fieldelement.NewFieldElement("0")
	b := fieldelement.NewFieldElement("7")
	alpha := x.Pow("3").Add(b)
	beta := alpha.Sqrt()
	var evenBeta fieldelement.FieldElement
	var oddBeta fieldelement.FieldElement
	if big.NewInt(0).Mod(beta.Number, big.NewInt(2)).Cmp(big.NewInt(0)) == 0 {
		evenBeta = beta
		Pbig, _ := big.NewInt(0).SetString(P, 16)
		oddBeta = fieldelement.NewFieldElement(big.NewInt(0).Sub(Pbig, beta.Number).String())
	} else {
		Pbig, _ := big.NewInt(0).SetString(P, 16)
		evenBeta = fieldelement.NewFieldElement(big.NewInt(0).Sub(Pbig, beta.Number).String())
		oddBeta = beta
	}
	if isEven {
		return &Point256{
			&x,
			&evenBeta,
			a,
			b,
		}
	}
	return &Point256{
		&x,
		&oddBeta,
		a,
		b,
	}

}

func (point256 *Point256) hash160(compressed bool) string {
	return utils.Hash160(point256.SEC(compressed))
}

//GetAddress will get the address from  the pubkey
func (point256 *Point256) GetAddress(compressed bool, testnet bool) string {
	hashedSEC := point256.hash160(compressed)
	var prefix string
	if testnet {
		prefix = "6f"
	} else {
		prefix = "00"
	}
	hashedWithPrefix := utils.Encode58CheckSum(prefix + hashedSEC)
	return hashedWithPrefix
}
