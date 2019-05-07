package core

import (
	"bytes"
	"github.com/StillFantastic/go-blockchain/wallet"
)

type TXInput struct {
	Txid			[]byte
	Vout			int
	Signature []byte
	PubKey		[]byte
}

type TXInputs struct {
	Inputs []TXInput
}

func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := wallet.HashPubKey(in.PubKey)
	return bytes.Compare(lockingHash, pubKeyHash) == 0
}


