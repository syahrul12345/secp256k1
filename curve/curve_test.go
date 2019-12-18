package curve

import (
	"fmt"
	"testing"
)

type coordinates map[string]string

func TestOnCurve(t *testing.T) {
	const prime int64 = 223
	const a uint64 = 0
	const b uint64 = 7
	validPoints := []coordinates{
		{
			"x": "192",
			"y": "105",
		},
		{
			"x": "17",
			"y": "56",
		},
		{
			"x": "1",
			"y": "193",
		},
	}
	for _, validPoint := range validPoints {
		// We create a new point
		_, notOnCurve := NewTestPoint(
			validPoint["x"],
			validPoint["y"],
			a,
			b,
			prime,
		)
		if notOnCurve != nil {
			t.Errorf("Expected point (%s,%s) to be on the curve, but it's not", validPoint["x"], validPoint["y"])
		}

	}
}

func TestNotOnCurve(t *testing.T) {
	const prime int64 = 223
	const a uint64 = 0
	const b uint64 = 7
	validPoints := []coordinates{
		{
			"x": "200",
			"y": "119",
		},
		{
			"x": "42",
			"y": "99",
		},
	}
	for _, validPoint := range validPoints {
		// We create a new point
		_, notOnCurve := NewTestPoint(
			validPoint["x"],
			validPoint["y"],
			a,
			b,
			prime,
		)
		if notOnCurve == nil {
			t.Errorf("Expected point (%s,%s) to be not on the curve, but it's on the curve", validPoint["x"], validPoint["y"])
		}

	}
}

func TestAdd(t *testing.T) {
	const prime int64 = 223
	const a uint64 = 0
	const b uint64 = 7
	var additions []coordinates = []coordinates{
		{
			"x1": "192",
			"y1": "105",
			"x2": "17",
			"y2": "56",
			"x3": "170",
			"y3": "142",
		},
		{
			"x1": "47",
			"y1": "71",
			"x2": "117",
			"y2": "141",
			"x3": "60",
			"y3": "139",
		},
		{
			"x1": "143",
			"y1": "98",
			"x2": "76",
			"y2": "66",
			"x3": "47",
			"y3": "71",
		},
	}
	for _, set := range additions {
		point1, _ := NewTestPoint(
			set["x1"],
			set["y1"],
			a,
			b,
			prime,
		)
		point2, _ := NewTestPoint(
			set["x2"],
			set["y2"],
			a,
			b,
			prime,
		)
		point3, _ := NewTestPoint(
			set["x3"],
			set["y3"],
			a,
			b,
			prime,
		)
		result, addError := point1.Add(point2)
		if addError != nil {
			t.Errorf("Failed to add Point(%v,%v,%v,%v) with Point(%v,%v,%v,%v) with prime: %d due to addition error",
				set["x1"],
				set["y1"],
				a,
				b,
				set["x2"],
				set["y2"],
				a,
				b,
				prime,
			)
		}
		expected := result.Equals(point3)
		if !expected {
			t.Errorf("The addition of Point(%v,%v,%v,%v) with Point(%v,%v,%v,%v) did not result in Point(%v,%v,%v,%v) ",
				set["x1"],
				set["y1"],
				a,
				b,
				set["x2"],
				set["y2"],
				a,
				b,
				set["x3"],
				set["y3"],
				a,
				b,
			)
		}
	}
}

func TestTrueDiv(t *testing.T) {
	const prime int64 = 223
	const a uint64 = 0
	const b uint64 = 7
	var validPoints []coordinates = []coordinates{
		{
			"coeff": "2",
			"x1":    "192",
			"y1":    "105",
			"x2":    "49",
			"y2":    "71",
		},
		{
			"coeff": "2",
			"x1":    "143",
			"y1":    "98",
			"x2":    "64",
			"y2":    "168",
		},
		{
			"coeff": "2",
			"x1":    "47",
			"y1":    "71",
			"x2":    "36",
			"y2":    "111",
		},
		{
			"coeff": "4",
			"x1":    "47",
			"y1":    "71",
			"x2":    "194",
			"y2":    "51",
		},
		{
			"coeff": "8",
			"x1":    "47",
			"y1":    "71",
			"x2":    "116",
			"y2":    "55",
		},
		{
			"coeff": "21",
			"x1":    "47",
			"y1":    "71",
			"x2":    "nil",
			"y2":    "nil",
		},
	}
	for _, set := range validPoints {
		point1, _ := NewTestPoint(
			set["x1"],
			set["y1"],
			a,
			b,
			prime,
		)
		point2, err2 := NewTestPoint(
			set["x2"],
			set["y2"],
			a,
			b,
			prime,
		)
		if err2 != nil {
			fmt.Printf("Point(%v,%v,%v,%v) doesnt exist on the curve. Aborting test without failing it.",
				set["x2"],
				set["y2"],
				a,
				b,
			)
			break
		}
		resultPoint, _ := point1.Mul(set["coeff"])
		equals := resultPoint.Equals(point2)
		if !equals {
			t.Errorf("Attempted the multiplacation of Point(%v,%v,%v,%v) with %v , and check that it gives Point(%v,%v,%v,%v) but it gave Point(%v,%v,%v,%v) instead",
				set["x1"],
				set["y1"],
				a,
				b,
				set["coeff"],
				set["x2"],
				set["y2"],
				a,
				b,
				resultPoint.X.Number,
				resultPoint.Y.Number,
				a,
				b,
			)
		}

	}
}
