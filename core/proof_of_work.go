package core

import (
	"math/big"
	"crypto/sha256"
	"github.com/StillFantastic/go-blockchain/tool"
	"bytes"
	"fmt"
)

// The difficulty of proof of work
// The bigger targetBits is, the harder to mine a block
const targetBits = 15; 

type ProofOfWork struct {
	block *Block
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	target.Lsh(target, 256 - targetBits)

	pow := &ProofOfWork{b, target}
	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			tool.IntToHex(pow.block.Timestamp),
			tool.IntToHex(int64(nonce)),
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			tool.IntToHex(int64(pow.block.Height)),
		},
		[]byte{},
	)
	return data
}

func (pow *ProofOfWork) Mine() (int, []byte) {
	var hashInt big.Int	
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block.\n")
	for {
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])

		if hashInt.Cmp(pow.target) == -1 {
			break
		} else {
			nonce++;
		}
	}
	
	fmt.Printf("%x\n", hash)
	fmt.Println()
	return nonce, hash[:]
}

func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1
	return isValid
}
