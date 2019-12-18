package main

import (
	"fmt"
	"secp256k1/creepto"
	"strconv"
)

func main() {
	//Create a new private key
	privKey := creepto.CreateNewPrivateKey()
	//Generate the pubKey from privateKey
	pubKey := privKey.PublicKey
	signature, message := privKey.Sign("hello world213124234123123")
	//Verifying signature
	res := pubKey.Verify(message, signature)
	//This will give true
	fmt.Println(res)
}

func intToHexadecimal(integer string) string {
	n, _ := strconv.ParseUint(integer, 16, 64)
	coefficientString := strconv.Itoa(int(n))
	return "0x" + coefficientString

}
