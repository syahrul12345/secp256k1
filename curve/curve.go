package curve

import (
	"fmt"
	"math/big"
	"secp256k1/fieldelement"

	"github.com/bitherhq/go-bither/common/hexutil"
)

type point struct {
	X *fieldelement.FieldElement
	Y *fieldelement.FieldElement
	A fieldelement.FieldElement
	B fieldelement.FieldElement
}

const (
// G is the generator point

)

type errorMessage struct {
	s string
}

type signature struct {
	R *fieldelement.FieldElement
	S *fieldelement.FieldElement
}

func (e *errorMessage) Error() string {
	return e.s
}

//NewPoint is a constructor function to create the new point
func NewPoint(x string, y string) (*point, error) {
	feX := fieldelement.NewFieldElement(x)
	feY := fieldelement.NewFieldElement(y)
	feA := fieldelement.NewFieldElement(hexutil.EncodeUint64(0))
	feB := fieldelement.NewFieldElement(hexutil.EncodeUint64(7))
	// Check if point exists on the curve
	onCurve := func(
		x fieldelement.FieldElement,
		y fieldelement.FieldElement,
		a fieldelement.FieldElement,
		b fieldelement.FieldElement) bool {
		left := y.Pow(2)
		r1 := x.Pow(3)
		r2 := x.Mul(a)
		r3 := b
		r4 := r1.Add(r2)
		right := r4.Add(r3)
		return left.Equals(right)
	}(fieldelement.NewFieldElement(x),
		fieldelement.NewFieldElement(y),
		fieldelement.NewFieldElement(hexutil.EncodeUint64(0)),
		fieldelement.NewFieldElement(hexutil.EncodeUint64(7)))

	if onCurve == false {
		return nil,
			&errorMessage{"The point doesnt exist on the curve"}
	}

	return &point{
		&feX,
		&feY,
		feA,
		feB,
	}, nil
}

//Helper Functions
//Equals will check if point1 is equals to point 2
func (point1 *point) Equals(point2 *point) bool {
	x1 := point1.X
	y1 := point1.Y
	a1 := point1.A
	b1 := point1.B
	x2 := point2.X
	y2 := point2.Y
	a2 := point2.A
	b2 := point2.B
	return x1.Equals(*x2) && y1.Equals(*y2) && a1.Equals(a2) && b1.Equals(b2)
}

//NotEquals will check if point1 is not equals to point2
func (point1 *point) NotEquals(point2 *point) bool {
	x1 := point1.X
	y1 := point1.Y
	a1 := point1.A
	b1 := point1.B
	x2 := point2.X
	y2 := point2.Y
	a2 := point2.A
	b2 := point2.B
	return x1.NotEquals(*x2) || y1.NotEquals(*y2) || a1.NotEquals(a2) || b1.NotEquals(b2)
}

//Multiplies the a point with a coefficient
func (point1 *point) Mul(coefficient string) (*point, error) {
	coeff, _ := hexutil.DecodeBig(coefficient)
	tempPoint := point1
	result := &point{
		nil,
		nil,
		point1.A,
		point1.B,
	}

	for coeff.And(coeff, big.NewInt(1)) != big.NewInt(1) {
		fmt.Println("adding")
		tempResult, addErr := result.Add(tempPoint)
		if addErr != nil {
			fmt.Println(addErr)
		}
		result = tempResult
		fmt.Println(result.X)
	}
	tempPoint.Add(tempPoint)
	coeff.Rsh(coeff, 1)
	fmt.Println("")
	return result, nil

}

//Adds point1 to point2
func (point1 *point) Add(point2 *point) (*point, error) {
	//Check if the points are on the same curve
	a1 := point1.A
	b1 := point1.B
	a2 := point2.A
	b2 := point2.B
	if a1.NotEquals(a2) || b1.NotEquals(b2) {
		return nil, &errorMessage{"They don't exist on the same point"}
	}

	x1 := point1.X
	y1 := point1.Y
	x2 := point2.X
	y2 := point2.Y
	//Case 0:
	if x1 == nil {
		return point2, nil
	}
	if x2 == nil {
		return point1, nil
	}
	// Case 1: Point @ Infinity. X is equals, but Y different. Vertical line.
	if x1.Equals(*x2) && y1.NotEquals(*y2) {
		return &point{
			nil,
			nil,
			a1,
			b2,
		}, &errorMessage{"Point of infinity"}
	}
	//Case 2: Point 1 and Point 2 are totally differnet
	if x1.NotEquals(*x2) && y1.NotEquals(*y2) {
		s := (y1.Sub(*y2)).TrueDiv(x1.Sub(*x2))
		x3 := s.Pow(2).Sub(*x1).Sub(*x2)
		y3 := s.Mul(x1.Sub(x3)).Sub(*y1)
		return &point{
			&x3,
			&y3,
			a1,
			b2,
		}, nil
	}
	//Case 3: The tangent of the point forms avertical line
	if point1.Equals(point2) && point1.Y.Equals(point1.X.Mul(fieldelement.NewFieldElement(hexutil.EncodeUint64(0)))) {
		return &point{
			nil,
			nil,
			a1,
			b2,
		}, &errorMessage{"Tagent forms a vertical line"}
	}
	//Case 4: The two points are exactly the same!
	if point1.Equals(point2) {
		s := x1.Pow(2).Add(x1.Pow(2).Add(x1.Pow(2))).Add(a1).TrueDiv(y1.Add(*y1))
		x3 := s.Pow(2).Sub(x1.Add(*x1))
		y3 := s.Mul(x1.Sub(x3)).Sub(*y1)
		return &point{
			&x3,
			&y3,
			a1,
			b1,
		}, nil
	}

	return nil, nil
}

func (point *point) Verify(zs string, sig *signature) {
	sInv := sig.S.Pow(-1)
	z := fieldelement.NewFieldElement(zs)
	u := z.Mul(sInv)
	v := sig.R.Mul(sInv)
	G, _ := NewPoint("0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", "0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8")
	first, _ := G.Mul(hexutil.EncodeBig(u.Number))
	fmt.Println(first)
	second, _ := point.Mul(hexutil.EncodeBig(v.Number))
	fmt.Println(second)
	// total, _ := first.Add(second)
	// print(total.X.Number)

}

//NewSignature creats a new signature object which stores the required variables
func NewSignature(R string, S string) *signature {
	r := fieldelement.NewFieldElement(R)
	s := fieldelement.NewFieldElement(S)
	return &signature{
		&r,
		&s,
	}
}
