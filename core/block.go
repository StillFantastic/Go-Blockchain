package core

import (
	"time"
	"encoding/gob"
	"bytes"
	"log"
)

type Block struct {
	Timestamp			int64
	Transactions  []*Transaction
	PrevBlockHash	[]byte
	Hash					[]byte
	Nonce					int
	Height				int
}

// Generate a new block
func NewBlock(transactions []*Transaction, prevBlockHash []byte, height int) *Block {
	block := &Block{time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0, height}
	pow := NewProofOfWork(block)
	nonce, hash := pow.Mine()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// This will serialize the block
func (b *Block) Serialize() []byte {
	var s bytes.Buffer
	encoder := gob.NewEncoder(&s)
	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}
	return s.Bytes()
}

// Deserialize a block
func DeserializeBlock(d []byte) *Block {
	var b Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&b)
	if err != nil {
		log.Panic(err)
	}
	return &b
}

// Generate a genesis block
func NewGenesisBlock(coinbase *Transaction) *Block {
	return NewBlock([]*Transaction{coinbase}, []byte{}, 1)
}

func (block *Block) HashTransactions() []byte {
	var transactions [][]byte

	for _, tx := range block.Transactions {
		transactions = append(transactions, tx.Serialize())
	}
	tree := NewMerkleTree(transactions)

	return tree.Root.Data
}
