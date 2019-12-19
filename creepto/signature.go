package creepto

import (
	"fmt"
)

//Signature represents a string
type Signature struct {
	R string
	S string
}

//NewSignature will create a new Signature
func NewSignature(r string, s string) *Signature {
	return &Signature{
		r,
		s,
	}
}

func (sig *Signature) DER() string {
	rbin := "0x37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6"
	sbin := "0x37206a0610995c58074999cb9767b87af4c4978db68c06e8e6e81d282047a7c6"
	totalbin := rbin + sbin
	b := []byte(totalbin)
	fmt.Println(len(b))
	return "l"
}
