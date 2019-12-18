package fieldelement

import (
	"fmt"
	"math"
	"math/big"
	"strconv"

	"github.com/bitherhq/go-bither/common/hexutil"
)

//These are the fixed values
var (
	// N float64  = math.Pow(2, 256) - math.Pow(2, 32) - 977
	first  *big.Int = big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0))
	second *big.Int = big.NewInt(0).Exp(big.NewInt(2), big.NewInt(32), big.NewInt(0))
	third  *big.Int = big.NewInt(977)
	n      *big.Int = big.NewInt(0).Sub(big.NewInt(0).Sub(first, second), third)
)

//FieldElement represetns a point in the elliptic curve
type FieldElement struct {
	Number *big.Int
	Prime  *big.Int
}

//NewFieldElement is a  constructor to create the fieldElement struct outside of the modesl pacakage
func NewFieldElement(ridiculouslyLargeNumber string) FieldElement {
	// This function only checks for the hex
	_, alreadyHex := strconv.ParseUint(ridiculouslyLargeNumber, 16, 64)
	if alreadyHex != nil {
		decodedNumber, decodeErr := hexutil.DecodeBig(ridiculouslyLargeNumber)
		if decodeErr != nil {
			//Handles if it's a point at infinity
			decodedNumber = big.NewInt(math.MaxInt64)
		}
		tempFieldElement := FieldElement{
			decodedNumber,
			n,
		}

		return tempFieldElement
	}
	str, _ := strconv.Atoi(ridiculouslyLargeNumber)
	formattedNumber := strconv.FormatInt(int64(str), 16)
	decodedNumber, decodeErr := hexutil.DecodeBig("0x" + formattedNumber)
	if decodeErr != nil {
		fmt.Println(decodeErr)
	}
	tempFieldElement := FieldElement{
		decodedNumber,
		n,
	}
	return tempFieldElement
}

//Creates a fieldElement for testing purposes
func NewTestingFieldElement(ridiculouslyLargeNumber string, testPrime int64) FieldElement {
	// This function only checks for the hex
	ntest := big.NewInt(testPrime)
	_, alreadyHex := strconv.ParseUint(ridiculouslyLargeNumber, 16, 64)
	if alreadyHex != nil {
		decodedNumber, decodeErr := hexutil.DecodeBig(ridiculouslyLargeNumber)
		if decodeErr != nil {
			decodedNumber = big.NewInt(math.MaxInt64)
		}
		tempFieldElement := FieldElement{
			decodedNumber,
			ntest,
		}

		return tempFieldElement
	}
	str, _ := strconv.Atoi(ridiculouslyLargeNumber)
	formattedNumber := strconv.FormatInt(int64(str), 16)
	decodedNumber, decodeErr := hexutil.DecodeBig("0x" + formattedNumber)
	if decodeErr != nil {
		fmt.Println(decodeErr)
	}
	tempFieldElement := FieldElement{
		decodedNumber,
		ntest,
	}

	return tempFieldElement
}

//Helper function to quickly convert if the int is in int and not hexadecimal
func intToHexadecimal(integer string) string {
	n, _ := strconv.ParseUint(integer, 16, 64)
	coefficientString := strconv.Itoa(int(n))
	return "0x" + coefficientString
}

//Equals check if the two fieldelements are equal
func (element1 FieldElement) Equals(element2 FieldElement) bool {
	numEquals := (element1.Number.Cmp(element2.Number))
	primeEquals := (element1.Prime.Cmp(element2.Prime))
	if numEquals == 0 && primeEquals == 0 {
		return true
	}
	return false
}

//NotEquals checks if two elements are not equal to each other
func (element1 FieldElement) NotEquals(element2 FieldElement) bool {
	numEquals := (element1.Number.Cmp(element2.Number))
	primeEquals := (element1.Prime.Cmp(element2.Prime))
	if numEquals != 0 || primeEquals != 0 {
		return true
	}
	return false
}

//Add to fieldelements together
func (element1 FieldElement) Add(element2 FieldElement) FieldElement {
	before := big.NewInt(0)
	before.Set(element1.Number)
	sum := before.Add(element1.Number, element2.Number)
	after := big.NewInt(0)
	after.Set(sum)
	num := after.Mod(sum, element1.Prime)
	return FieldElement{num, element1.Prime}
}

//Sub will subtract element2 from elemet1
func (element1 FieldElement) Sub(element2 FieldElement) FieldElement {
	before := big.NewInt(0)
	before.Set(element1.Number)
	sum := before.Sub(element1.Number, element2.Number)
	after := big.NewInt(0)
	after.Set(sum)
	num := after.Mod(sum, element1.Prime)
	return FieldElement{num, element1.Prime}
}

//Mul will subtract elememnt2 from element 1
func (element1 FieldElement) Mul(element2 FieldElement) FieldElement {
	before := big.NewInt(0)
	before.Set(element1.Number)
	sum := before.Mul(element1.Number, element2.Number)
	after := big.NewInt(0)
	after.Set(sum)
	num := after.Mod(sum, element1.Prime)
	return FieldElement{num, element1.Prime}
}

//Pow will exponnent the fieldl  e ment yegv power
func (element1 FieldElement) Pow(exponent int) FieldElement {
	exp := big.NewInt(int64(exponent))
	one := big.NewInt(1)
	primeLess := big.NewInt(0).Sub(element1.Prime, one)
	n := exp.Mod(exp, primeLess)
	num := exp.Exp(element1.Number, n, element1.Prime)
	return FieldElement{num, element1.Prime}
}

//Makes a point truely divisible
func (element1 FieldElement) TrueDiv(element2 FieldElement) FieldElement {
	var primeLess *big.Int
	two := big.NewInt(2)
	primeLess = big.NewInt(0).Sub(element1.Prime, two)
	divisor := big.NewInt(0).Exp(element2.Number, primeLess, element1.Prime)
	num := big.NewInt(0).Mod(element1.Number.Mul(element1.Number, divisor), element1.Prime)
	return FieldElement{num, element1.Prime}
}
