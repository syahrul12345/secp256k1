package curve

import (
	"fmt"
	"math/big"
	"secp256k1/fieldelement"
	"secp256k1/utils"
)

type Point struct {
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

func (e *errorMessage) Error() string {
	return e.s
}

//NewPoint is a constructor function to create the new point
func NewPoint(x string, y string) (*Point, error) {

	var feX *fieldelement.FieldElement
	var feY *fieldelement.FieldElement
	feA := fieldelement.NewFieldElement(utils.FromBigInt(big.NewInt(0)))
	feB := fieldelement.NewFieldElement(utils.FromBigInt(big.NewInt(7)))
	if x == "nil" {
		feX = nil
	}
	if y == "nil" {
		feY = nil
	} else {
		feXVal := fieldelement.NewFieldElement(x)
		feYVal := fieldelement.NewFieldElement(y)
		feX = &feXVal
		feY = &feYVal
		// Check if point exists on the curve. SKip if the X or y == nil
		onCurve := func(
			x fieldelement.FieldElement,
			y fieldelement.FieldElement,
			a fieldelement.FieldElement,
			b fieldelement.FieldElement) bool {
			left := y.Pow("2")
			r1 := x.Pow("3")
			r2 := x.Mul(a)
			r3 := b
			r4 := r1.Add(r2)
			right := r4.Add(r3)
			return left.Equals(right)
		}(fieldelement.NewFieldElement(x),
			fieldelement.NewFieldElement(y),
			fieldelement.NewFieldElement(utils.FromBigInt(big.NewInt(0))),
			fieldelement.NewFieldElement(utils.FromBigInt(big.NewInt(7))))
		if onCurve == false {
			return nil,
				&errorMessage{"The point doesnt exist on the curve"}
		}
	}
	return &Point{
		feX,
		feY,
		feA,
		feB,
	}, nil
}

//Helper Functions
//Equals will check if point1 is equals to point 2
func (point1 *Point) Equals(point2 *Point) bool {
	x1 := point1.X
	y1 := point1.Y
	a1 := point1.A
	b1 := point1.B
	x2 := point2.X
	y2 := point2.Y
	a2 := point2.A
	b2 := point2.B
	if x1 == nil || x2 == nil || y1 == nil || y2 == nil {
		if x1 == x2 && y1 == y2 {
			return true
		}
		return false
	}
	return x1.Equals(*x2) && y1.Equals(*y2) && a1.Equals(a2) && b1.Equals(b2)
}

//NotEquals will check if point1 is not equals to point2
func (point1 *Point) NotEquals(point2 *Point) bool {
	x1 := point1.X
	y1 := point1.Y
	a1 := point1.A
	b1 := point1.B
	x2 := point2.X
	y2 := point2.Y
	a2 := point2.A
	b2 := point2.B
	if x1 == nil || x2 == nil || y1 == nil || y2 == nil {
		if x1 == x2 && y1 == y2 {
			return false
		}
		return true
	}
	return x1.NotEquals(*x2) || y1.NotEquals(*y2) || a1.NotEquals(a2) || b1.NotEquals(b2)
}

//Mul will multiply the point with a string
func (point1 *Point) Mul(coefficient string) (*Point, error) {
	coeff := utils.ToBigInt(coefficient)
	fmt.Println(coeff)
	current := point1
	result := &Point{
		nil,
		nil,
		point1.A,
		point1.B,
	}
	//Binary expansion
	for coeff.Cmp(big.NewInt(0)) > 0 {
		// keep adding to ther result if the rightmost bit is 1
		if big.NewInt(0).And(coeff, big.NewInt(1)).Int64() == 1 {
			result, _ = result.Add(current)
		}
		current, _ = current.Add(current)
		coeff.Rsh(coeff, 1)
	}
	return result, nil
}

//Adds point1 to point2
func (point1 *Point) Add(point2 *Point) (*Point, error) {
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
		return &Point{
			nil,
			nil,
			a1,
			b2,
		}, &errorMessage{"Point of infinity"}
	}
	//Case 2: Point 1 and Point 2 are totally differnet
	if x1.NotEquals(*x2) && y1.NotEquals(*y2) {
		numerator := y1.Sub(*y2)
		denominator := x1.Sub(*x2)
		s := numerator.TrueDiv(denominator)
		x3 := s.Pow("2").Sub(*x1).Sub(*x2)
		y3 := s.Mul(x1.Sub(x3)).Sub(*y1)
		return &Point{
			&x3,
			&y3,
			a1,
			b2,
		}, nil
	}
	//Case 3: The tangent of the point forms avertical line
	if point1.Equals(point2) && point1.Y.Equals(point1.X.Mul(fieldelement.NewFieldElement(utils.FromBigInt(big.NewInt(0))))) {
		return &Point{
			nil,
			nil,
			a1,
			b2,
		}, &errorMessage{"Tagent forms a vertical line"}
	}
	//Case 4: The two points are exactly the same!
	if point1.Equals(point2) {
		s := x1.Pow("2").Add(x1.Pow("2").Add(x1.Pow("2"))).Add(a1).TrueDiv(y1.Add(*y1))
		x3 := s.Pow("2").Sub(x1.Add(*x1))
		y3 := s.Mul(x1.Sub(x3)).Sub(*y1)
		return &Point{
			&x3,
			&y3,
			a1,
			b1,
		}, nil
	}

	return nil, nil
}
