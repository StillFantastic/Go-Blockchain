package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"log"
	"bytes"
	"golang.org/x/crypto/ripemd160"
	"github.com/StillFantastic/go-blockchain/tool"
)

const version = byte(0x00)
const addressChecksumLen = 4

type Wallet struct {
	PrivateKey ecdsa.PrivateKey
	PublicKey []byte
}

func NewWallet() *Wallet {
	private, public := newKeyPair()
	return &Wallet{private, public}
}

func newKeyPair() (ecdsa.PrivateKey, []byte) {
	curve := elliptic.P256()
	private, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		log.Panic(err)
	}
	publicKey := append(private.PublicKey.X.Bytes(), private.PublicKey.Y.Bytes()...)
	return *private, publicKey
}

func (w Wallet) GetAddress() []byte {
	pubKeyHash := HashPubKey(w.PublicKey)
	
	versionedPayload := append([]byte{version}, pubKeyHash...)
	checksum := checksum(versionedPayload)

	fullPayload := append(versionedPayload, checksum...)
	address := tool.Base58Encode(fullPayload)

	return address
}

func HashPubKey(pubKey []byte) []byte {
	pubKeySHA256 := sha256.Sum256(pubKey)

	RIPEMD160Hasher := ripemd160.New()
	_, err := RIPEMD160Hasher.Write(pubKeySHA256[:])
	if err != nil {
		log.Panic(err)
	}
	pubKeyRIPEMD160 := RIPEMD160Hasher.Sum(nil)

	return pubKeyRIPEMD160
}

func checksum(payload []byte) []byte {
	firstHash := sha256.Sum256(payload)
	secondHash := sha256.Sum256(firstHash[:])

	return secondHash[:addressChecksumLen]
}

func ValidateAddress(address string) bool {
	pubKeyHash := tool.Base58Decode([]byte(address))
	myChecksum := pubKeyHash[len(pubKeyHash) - addressChecksumLen:]
	targetChecksum := checksum(pubKeyHash[0:len(pubKeyHash) - addressChecksumLen])

	return bytes.Compare(myChecksum, targetChecksum) == 0
}
