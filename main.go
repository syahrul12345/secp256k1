package main

import "secp256k1/curve"

func main() {
	// fieldElement2 := fieldelement.NewFieldElement(11)
	// fmt.Println(fieldElement.Pow(-4).Mul(fieldElement2))
	// G, _ := curve.NewPoint("0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", "0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8")
	// point1, _ := curve.NewPoint("0x5cbdf0646e5db4eaa398f365f2ea7a0e3d419b7e0330e39ce92bddedcac4f9bc", "0x6aebca40ba255960a3178d6d861a54dba813d0b813fde7b5a5082628087264da")
	// n, _ := hexutil.DecodeUint64("0xfffffffffffffffffffffffffffffffebaaedce6af48a03bbfd25e8cd0364141")
	// fmt.Println(*G.X)
	// fmt.Println(*point1.X)
	// fmt.Println(n)
	point, _ := curve.NewPoint("0x887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c", "0x61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34")
	z := "0xec208baa0fc1c19f708a9ca96fdeff3ac3f230bb4a7ba4aede4942ad003c0f60"
	signature := curve.NewSignature("0xac8d1c87e51d0d441be8b3dd5b05c8795b48875dffe00b7ffcfac23010d3a395", "0x68342ceff8935ededd102dd876ffd6ba72d6a427a3edb13d26eb0781cb423c4")
	point.Verify(z, signature)
}
