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
	signature, message := privKey.Sign("hello wo123213123rld213124234123123")
	//Verifying signature
	status := pubKey.Verify(message, signature)
	fmt.Println(status)
	//This will give true
	secPub := pubKey.SEC(true)
	fmt.Println(secPub)
	secPub = pubKey.SEC(false)
	fmt.Println(secPub)
	fmt.Println("---------------------------------------------------------------------------")
}

func intToHexadecimal(integer string) string {
	n, _ := strconv.ParseUint(integer, 16, 64)
	coefficientString := strconv.Itoa(int(n))
	return "0x" + coefficientString

}
