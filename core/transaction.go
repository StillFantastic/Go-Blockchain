package core

import (
	"bytes"
	"crypto/sha256"
	"crypto/rand"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/gob"
	"math/big"
	"log"
	"fmt"
	"time"
	"strconv"
	"encoding/hex"
	"github.com/StillFantastic/go-blockchain/wallet"
)

const subsidy = 10

type Transaction struct {
	ID		[]byte
	Vin		[]TXInput
	Vout	[]TXOutput
}

// Serialize transaction
func (transaction Transaction) Serialize() []byte {
	var b bytes.Buffer
	encoder := gob.NewEncoder(&b)
	err := encoder.Encode(transaction)
	if err != nil {
		log.Panic(err)
	}
	return b.Bytes()
}

// Set ID for the transaction
func (transaction *Transaction) SetID() {
	var hash [32]byte
	hash = sha256.Sum256(transaction.Serialize())
	transaction.ID = hash[:]
}

// Generate a coinbase transaction
func NewCoinbaseTX(to, data string) *Transaction{
	if data == "" {
		data = fmt.Sprintf("Reward to %s, at %s", to, strconv.FormatInt(time.Now().Unix(), 10))
	}

	txin := TXInput{[]byte{}, -1, nil, []byte(data)}
	txout := NewTXOutput(subsidy, to)
	tx := Transaction{nil, []TXInput{txin}, []TXOutput{*txout}}
	tx.SetID()

	return &tx
}

// Check if the transaction is a coinbase transaction
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Vin) == 1 && len(tx.Vin[0].Txid) == 0 && tx.Vin[0].Vout == -1
}

// Generate a new transaction with unspent transactions
func (bc *Blockchain) NewUTXOTransaction(wt *wallet.Wallet, to string, amount int, UTXOSet *UTXOSet) *Transaction {
	var inputs []TXInput
	var outputs []TXOutput

	pubKeyHash := wallet.HashPubKey(wt.PublicKey)
	acc, unspentOutputs := UTXOSet.FindSpendableOutputs(pubKeyHash, amount)

	if acc < amount {
		log.Panic("ERROR: Not enough funds.")
	}
	
	for txid, outputs := range unspentOutputs {
		txID, err := hex.DecodeString(txid)
		if err != nil {
			log.Panic(err)
		}
		for _, outputID := range outputs {
			inputs = append(inputs, TXInput{txID, outputID, nil, wt.PublicKey})
		}
	}

	from := string(wt.GetAddress())
	outputs = append(outputs, *NewTXOutput(amount, to))
	if acc > amount {
		outputs = append(outputs, *NewTXOutput(acc - amount, from))
	}
	tx := &Transaction{nil, inputs, outputs}
	tx.SetID()
	bc.SignTransaction(tx, wt.PrivateKey)
	return tx	
}

// Replicate a transaction only with some necessary information
func (tx Transaction) TrimmedCopy() Transaction{
	var inputs []TXInput
	var outputs []TXOutput

	for _, out := range tx.Vout {
		outputs = append(outputs, TXOutput{out.Value, out.PubKeyHash})
	}

	for _, in := range tx.Vin {
		inputs = append(inputs, TXInput{in.Txid, in.Vout, nil, nil})
	}

	return Transaction{nil, inputs, outputs}
}

// Sign a transaction 
func (tx *Transaction) Sign(privateKey ecdsa.PrivateKey, prevTXs map[string]Transaction) {
	if tx.IsCoinbase() {
		return
	}

	txCopy := tx.TrimmedCopy()

	for inID, vin := range txCopy.Vin {
		prevTX := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].PubKey = prevTX.Vout[vin.Vout].PubKeyHash
		txCopy.SetID()
		txHash := txCopy.ID
		txCopy.Vin[inID].PubKey = nil
		txCopy.ID = nil

		r, s, err := ecdsa.Sign(rand.Reader, &privateKey, txHash)
		if err != nil {
			log.Panic(err)
		}
		signature := append(r.Bytes(), s.Bytes()...)
		tx.Vin[inID].Signature = signature
	}
}

// Verify signatures of transaction inputs
func (tx *Transaction) Verify(prevTXs map[string]Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	txCopy := tx.TrimmedCopy()
	curve := elliptic.P256()	

	for inID, vin := range tx.Vin {
		prevTX := prevTXs[hex.EncodeToString(vin.Txid)]
		txCopy.Vin[inID].PubKey = prevTX.Vout[vin.Vout].PubKeyHash
		txHash := txCopy.ID
		txCopy.Vin[inID].PubKey = nil
		txCopy.ID = nil	
	
		r := big.Int{}
		s := big.Int{}
		sigLen := len(vin.Signature)
		r.SetBytes(vin.Signature[:(sigLen / 2)])
		s.SetBytes(vin.Signature[(sigLen / 2):])

		x := big.Int{}
		y := big.Int{}
		keyLen := len(vin.PubKey)
		x.SetBytes(vin.PubKey[:(keyLen / 2)])
		y.SetBytes(vin.PubKey[(keyLen / 2):])

		pubKey := ecdsa.PublicKey{curve, &x, &y}
		if ecdsa.Verify(&pubKey, txHash, &r, &s) == false {
			return false
		}
	}
	return true
}

func DeserializeTransaction(data []byte) Transaction {
	var transaction Transaction

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&transaction)
	if err != nil {
		log.Panic(err)
	}

	return transaction
}
