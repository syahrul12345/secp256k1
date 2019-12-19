package fieldelement

import (
	"testing"
)

func TestEquals(t *testing.T) {
	a := NewTestingFieldElement("2", 31)
	b := NewTestingFieldElement("2", 31)
	c := NewTestingFieldElement("0x5cbdf0646e5db4eaa398f365f2ea7a0e3d419b7e0330e39ce92bddedcac4f9bc", 31)
	d := NewTestingFieldElement("0x5cbdf0646e5db4eaa398f365f2ea7a0e3d419b7e0330e39ce92bddedcac4f9bc", 31)
	aEqualB := a.Equals(b)
	cEqualD := c.Equals(d)
	aEqualD := a.Equals(d)
	if !aEqualB {
		t.Errorf("a.Equals(b) = %t; want true", aEqualB)
	}
	if !cEqualD {
		t.Errorf("c.Equals(D) = %t; want true", cEqualD)
	}
	if aEqualD {
		t.Errorf("a.Equals(D) = %t; want false", aEqualD)
	}

}

func TestNotEquals(t *testing.T) {
	a := NewTestingFieldElement("2", 31)
	b := NewTestingFieldElement("2", 31)
	c := NewTestingFieldElement("0x5cbdf0646e5db4eaa398f365f2ea7a0e3d419b7e0330e39ce92bddedcac4f9bc", 31)
	d := NewTestingFieldElement("0x5cbdf0646e5db4eaa398f365f2ea7a0e3d419b7e0330e39ce92bddedcac4f9bc", 31)
	aNotEqualB := a.NotEquals(b)
	cNotEqualD := c.NotEquals(d)
	aNotEqualD := a.NotEquals(d)
	bNotEqualC := b.NotEquals(c)
	if aNotEqualB {
		t.Errorf("a.NotEquals(b) = %t, want false", aNotEqualB)
	}
	if cNotEqualD {
		t.Errorf("c.NotEquals(d) = %t, want false", cNotEqualD)
	}
	if !aNotEqualD {
		t.Errorf("a.NotEquals(d) = %t, want true", aNotEqualD)
	}
	if !bNotEqualC {
		t.Errorf("b.NotEquals(c) = %t, want true", bNotEqualC)
	}
}

func TestAdd(t *testing.T) {
	a := NewTestingFieldElement("2", 31)
	b := NewTestingFieldElement("15", 31)
	result1 := a.Add(b)
	expected1 := NewTestingFieldElement("17", 31)
	first := result1.Equals(expected1)
	if !first {
		t.Errorf("Field element addition of FieldElement(2) and FieldElement(15) gave FieldElement(%s)s , want FieldElement(17)", result1.Number)
	}
	c := NewTestingFieldElement("17", 31)
	d := NewTestingFieldElement("21", 31)
	result2 := c.Add(d)
	expected2 := NewTestingFieldElement("7", 31)
	second := result2.Equals(expected2)
	if !second {
		t.Errorf("Field element addition of FieldElement(17) and FieldElement(21) gave FieldElement(%s) ,want FieldElement(7)", result2.Number)
	}
}

func TestSub(t *testing.T) {
	a := NewTestingFieldElement("29", 31)
	b := NewTestingFieldElement("4", 31)
	result1 := a.Sub(b)
	expected1 := NewTestingFieldElement("25", 31)
	first := result1.Equals(expected1)
	if !first {
		t.Errorf("Field element subtraction of FieldElement(29) - FieldElement(4) gave FieldElement(%s), want FieldElement(25) ", result1.Number)
	}
	c := NewTestingFieldElement("15", 31)
	d := NewTestingFieldElement("30", 31)
	result2 := c.Sub(d)
	expected2 := NewTestingFieldElement("16", 31)
	second := result2.Equals(expected2)
	if !second {
		t.Errorf("Field element addition of FieldElement(15) and FieldElement(30) gave FieldElement(%s) ,want FieldElement(16)", result2.Number)
	}
}

func TestMul(t *testing.T) {
	a := NewTestingFieldElement("24", 31)
	b := NewTestingFieldElement("19", 31)
	result1 := a.Mul(b)
	expected1 := NewTestingFieldElement("22", 31)
	first := result1.Equals(expected1)
	if !first {
		t.Errorf("Field element addition of FieldElement(24) and FieldElement(19) gave FieldElement(%s) ,want FieldElement(16)", result1.Number)
	}
}

func TestPow(t *testing.T) {
	a := NewTestingFieldElement("17", 31)
	result1 := a.Pow("3")
	expected1 := NewTestingFieldElement("15", 31)
	first := result1.Equals(expected1)
	if !first {
		t.Errorf("Field element exponention of FieldElement(17) by Exponent 3 gave FieldElement(%s), want FieldElement(15)", result1.Number)
	}
	b := NewTestingFieldElement("5", 31)
	c := NewTestingFieldElement("18", 31)
	result2 := b.Pow("5").Mul(c)
	expected2 := NewTestingFieldElement("16", 31)
	second := result2.Equals(expected2)
	if !second {
		t.Errorf("Field element exponent of FieldElement(5) by Exponent 5, followed by Multiplication with FieldElement(18) gave FieldElement(%s),want FieldElement(16)", result2.Number)
	}
	result3 := a.Pow("-3")
	expected3 := NewTestingFieldElement("29", 31)
	third := result3.Equals(expected3)
	if !third {
		t.Errorf("Field element exponent of FieldELement(17) by Exponent -3, gave FieldElement(%s),want FieldElement(29)", result3.Number)
	}
}

func TestDiv(t *testing.T) {
	a := NewTestingFieldElement("3", 31)
	b := NewTestingFieldElement("24", 31)
	result1 := a.TrueDiv(b)
	expected1 := NewTestingFieldElement("4", 31)
	first := result1.Equals(expected1)
	if !first {
		t.Errorf("Field element divition of FieldElement(3)/FieldElement(24) gave FieldElement(%s), want FieldElement(4) ", result1.Number)
	}
	a = NewTestingFieldElement("17", 31)
	result2 := a.Pow("-3")
	expected2 := NewTestingFieldElement("29", 31)
	second := result2.Equals(expected2)
	if !second {
		t.Errorf("Field element exponent of FieldELement(17) by Exponent -3, gave FieldElement(%s),want FieldElement(29)", result2.Number)
	}
	c := NewTestingFieldElement("4", 31)
	d := NewTestingFieldElement("11", 31)
	result3 := c.Pow("-4").Mul(d)

	expected3 := NewTestingFieldElement("13", 31)
	third := result3.Equals(expected3)
	if !third {
		t.Errorf("Field element exponent of FieldElement(4) by Exponent -4 then multiply By FieldElement(11) gave FieldElement(%s), want FieldElement(13)", result3.Number)
	}
}
