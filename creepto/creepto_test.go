package creepto

import (
	"fmt"
	"math/big"
	"secp256k1/curve"
	"secp256k1/utils"
	"testing"

	"github.com/bitherhq/go-bither/common/hexutil"
)

type coordinates map[string]string

func TestPubPoint(t *testing.T) {
	// 2**128
	set1 := hexutil.EncodeBig(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(128), big.NewInt(0)))
	// 2**240 + 2**31
	set2 := hexutil.EncodeBig(big.NewInt(0).Add(big.NewInt(0).Exp(big.NewInt(2), big.NewInt(240), big.NewInt(0)), big.NewInt(0).Exp(big.NewInt(2), big.NewInt(31), big.NewInt(0))))
	var sets []coordinates = []coordinates{
		{
			"secret": "7",
			"x":      "0x5cbdf0646e5db4eaa398f365f2ea7a0e3d419b7e0330e39ce92bddedcac4f9bc",
			"y":      "0x6aebca40ba255960a3178d6d861a54dba813d0b813fde7b5a5082628087264da",
		},
		{
			"secret": "1485",
			"x":      "0xc982196a7466fbbbb0e27a940b6af926c1a74d5ad07128c82824a11b5398afda",
			"y":      "0x7a91f9eae64438afb9ce6448a1c133db2d8fb9254e4546b6f001637d50901f55",
		},
		{
			"secret": set1,
			"x":      "0x8f68b9d2f63b5f339239c1ad981f162ee88c5678723ea3351b7b444c9ec4c0da",
			"y":      "0x662a9f2dba063986de1d90c2b6be215dbbea2cfe95510bfdf23cbf79501fff82",
		},
		{
			"secret": set2,
			"x":      "0x9577ff57c8234558f293df502ca4f09cbc65a6572c842b39b366f21717945116",
			"y":      "0x10b49c67fa9365ad7b90dab070be339a1daf9052373ec30ffae4f72d5e66d053",
		},
	}
	for _, set := range sets {
		G, _ := curve.NewPoint("0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", "0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8")
		point := New256Point(set["x"], set["y"])
		result, _ := G.Mul(set["secret"])
		final := point.X.Number.Cmp(result.X.Number)
		if final != 0 {
			t.Errorf("Multyplying the secret number %v with G, the generator point should give the public key at Point(%v,%v), but got Point(%v,%v)",
				set["secret"],
				set["x"],
				set["y"],
				hexutil.EncodeBig(result.X.Number),
				hexutil.EncodeBig(result.Y.Number),
			)
		}

	}
}

func TestVerify(t *testing.T) {
	point := New256Point("0x887387e452b8eacc4acfde10d9aaf7f6d9a0f975aabb10d006e4da568744d06c", "0x61de6d95231cd89026e286df3b6ae4a894a3378e393e93a0f45b666329a0ae34")
	var sets []coordinates = []coordinates{
		{
			"z": "0xec208baa0fc1c19f708a9ca96fdeff3ac3f230bb4a7ba4aede4942ad003c0f60",
			"r": "0xac8d1c87e51d0d441be8b3dd5b05c8795b48875dffe00b7ffcfac23010d3a395",
			"s": "0x68342ceff8935ededd102dd876ffd6ba72d6a427a3edb13d26eb0781cb423c4",
		},
		{
			"z": "0x7c076ff316692a3d7eb3c3bb0f8b1488cf72e1afcd929e29307032997a838a3d",
			"r": "0xeff69ef2b1bd93a66ed5219add4fb51e11a840f404876325a1e8ffe0529a2c",
			"s": "0xc7207fee197d27c618aea621406f6bf5ef6fca38681d82b2f06fddbdce6feab6",
		},
	}
	for _, set := range sets {
		signature := NewSignature(set["r"], set["s"])
		verified := point.Verify(set["z"], signature)
		if !verified {
			t.Errorf("Verification of message Z(%s) with signatures (%s,%s) failed", set["z"], set["r"], set["s"])
		}
	}
}

func TestSec(t *testing.T) {
	set1 := hexutil.EncodeBig(big.NewInt(0).Exp(big.NewInt(999), big.NewInt(3), big.NewInt(0)))
	var sets []coordinates = []coordinates{
		{
			"coefficient":  set1,
			"compressed":   "039d5ca49670cbe4c3bfa84c96a8c87df086c6ea6a24ba6b809c9de234496808d5",
			"uncompressed": "049d5ca49670cbe4c3bfa84c96a8c87df086c6ea6a24ba6b809c9de234496808d56fa15cc7f3d38cda98dee2419f415b7513dde1301f8643cd9245aea7f3f911f9",
		},
		{
			"coefficient":  "123",
			"compressed":   "03a598a8030da6d86c6bc7f2f5144ea549d28211ea58faa70ebf4c1e665c1fe9b5",
			"uncompressed": "04a598a8030da6d86c6bc7f2f5144ea549d28211ea58faa70ebf4c1e665c1fe9b5204b5d6f84822c307e4b4a7140737aec23fc63b65b35f86a10026dbd2d864e6b",
		},
		{
			"coefficient":  "42424242",
			"compressed":   "02aee2e7d843f7430097859e2bc603abcc3274ff8169c1a469fee0f20614066f8e",
			"uncompressed": "04aee2e7d843f7430097859e2bc603abcc3274ff8169c1a469fee0f20614066f8e21ec53f40efac47ac1c5211b2123527e0e9b57ede790c4da1e72c91fb7da54a3",
		},
	}
	for _, set := range sets {
		G, _ := curve.NewPoint("0x79be667ef9dcbbac55a06295ce870b07029bfcdb2dce28d959f2815b16f81798", "0x483ada7726a3c4655da4fbfc0e1108a8fd17b448a68554199c47d08ffb10d4b8")
		//Generate the PublicKey
		point, err := G.Mul(set["coefficient"])
		if err != nil {
			fmt.Println(nil)
		}
		//Convert point to a point 256
		point256 := NormalTo256(point)
		compressed := point256.SEC(true)
		uncompressed := point256.SEC(false)
		if compressed != set["compressed"] {
			t.Errorf("Generated publicKey with coefficient %v multiplied by G, has an compressed SEC of %s , but want %s",
				set["coefficient"],
				compressed,
				set["compressed"],
			)
		}
		if uncompressed != set["uncompressed"] {
			t.Errorf("Generated publicKey with coefficient %v multiplied by G, has an uncompressed SEC of %s , but want %s",
				set["coefficient"],
				uncompressed,
				set["uncompressed"],
			)
		}

	}
}
func TestDer(t *testing.T) {
	// Generate 4 random hexadecimal numbers for r1,s1,r2,s2
	// Hexadecimals prepend 0x otherwise utils.ToBigInt() wont work
	r1 := utils.GenerateSecret()
	s1 := utils.GenerateSecret()
	r2 := utils.GenerateSecret()
	s2 := utils.GenerateSecret()
	// s := utils.GenerateSecret()
	var sets []coordinates = []coordinates{
		{
			"r": r1,
			"s": s1,
		},
		{
			"r": r2,
			"s": s2,
		},
	}
	for _, set := range sets {
		//Generate the signature
		sig := &Signature{
			utils.ToBigInt("0x" + set["r"]),
			utils.ToBigInt("0x" + set["s"]),
		}
		der := sig.DER()
		resSig := ParseDER(der)
		res := sig.Equals(resSig)
		if !res {
			t.Errorf("Test to parse the bin failed!")
		}
	}
}
