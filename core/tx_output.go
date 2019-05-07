package core

import (
	"github.com/StillFantastic/go-blockchain/tool"
	"bytes"
	"encoding/gob"
	"log"
)

type TXOutput struct {
	Value				int
	PubKeyHash	[]byte
}

type TXOutputs struct {
	Outputs []TXOutput
}

func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := tool.Base58Decode(address)
	pubKeyHash = pubKeyHash[1:len(pubKeyHash) - 4]
	out.PubKeyHash = pubKeyHash
}

func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

func NewTXOutput(amount int, address string) *TXOutput {
	txOutput := TXOutput{amount, nil}
	txOutput.Lock([]byte(address))

	return &txOutput
}

func (outs TXOutputs) Serialize() []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(outs)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func DeserializeOutputs(b []byte) TXOutputs {
	var outputs TXOutputs

	dec := gob.NewDecoder(bytes.NewReader(b))
	err := dec.Decode(&outputs)
	if err != nil {
		log.Panic(err)
	}
	
	return outputs
}
