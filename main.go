package main

import (
	"fmt"
	"secp256k1/creepto"
	"strconv"
)

func main() {
	// Create a new private key
	privKey := creepto.CreateNewPrivateKey()
	//Generate the pubKey from privateKey
	pubKey := privKey.PublicKey
	addr := pubKey.GetAddress()
	fmt.Println(addr)
}

func intToHexadecimal(integer string) string {
	n, _ := strconv.ParseUint(integer, 16, 64)
	coefficientString := strconv.Itoa(int(n))
	return "0x" + coefficientString

}
