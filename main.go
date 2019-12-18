package main

import (
	"secp256k1/crypto"
	"strconv"
)

func main() {
	G := crypto.New256Point(
		"0x887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c",
		"0x61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34")
	z := "0xec208baa0fc1c19f708a9ca96fdeff3ac3f230bb4a7ba4aede4942ad003c0f60"
	sig := crypto.NewSignature("0xac8d1c87e51d0d441be8b3dd5b05c8795b48875dffe00b7ffcfac23010d3a395", "0x68342ceff8935ededd102dd876ffd6ba72d6a427a3edb13d26eb0781cb423c4")
	G.Verify(z, sig)

}

func intToHexadecimal(integer string) string {
	n, _ := strconv.ParseUint(integer, 16, 64)
	coefficientString := strconv.Itoa(int(n))
	return "0x" + coefficientString

}
