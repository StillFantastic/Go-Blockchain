package tool

import (
	"bytes"
	"math/big"
)

var base58Alphabet = []byte("123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz")

func Base58Encode(b []byte) []byte {
	var base58 []byte

	now := big.NewInt(0).SetBytes(b)

	base := big.NewInt(int64(len(base58Alphabet)))
	zero := big.NewInt(0)
	mod := &big.Int{}

	for now.Cmp(zero) != 0 {
		now.DivMod(now, base, mod)
		base58 = append(base58, base58Alphabet[mod.Int64()])
	}

	for i := range b {
		if b[i] == 0x00 {
			base58 = append(base58, base58Alphabet[0])
		} else {
			break
		}
	}
	reverseBytes(base58)
	return base58
}

func reverseBytes(b []byte) {
	for i := 0; i < len(b) / 2; i++ {
		b[i], b[len(b) - 1 - i] = b[len(b) - i - 1], b[i]
	}
}

func Base58Decode(b []byte) []byte {
	ret := big.NewInt(0)
	zero := 0

	for i := range b {
		if b[i] == base58Alphabet[0] {
			zero++
		} else {
			break
		}
	}

	tmp := b[zero:]
	for _, x := range tmp {
		index := bytes.IndexByte(base58Alphabet, x)
		ret.Mul(ret, big.NewInt(58))
		ret.Add(ret, big.NewInt(int64(index)))
	}

	decoded := ret.Bytes()
	decoded = append(bytes.Repeat([]byte{byte(0x00)}, zero), decoded...)
	return decoded
}
