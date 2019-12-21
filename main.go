package main

import (
	"fmt"
	"secp256k1/creepto"
	"strconv"
)

func main() {
	// Create a new private key
	// pk := creepto.CreateNewPrivateKey()
	// fmt.Println(pk.PublicKey.GetAddress())
	pubKey := creepto.GetPublicKey("5000")
	fmt.Println(pubKey.GetAddress(false, false))

}

func intToHexadecimal(integer string) string {
	n, _ := strconv.ParseUint(integer, 16, 64)
	coefficientString := strconv.Itoa(int(n))
	return "0x" + coefficientString

}
