package core

import (
	"fmt"
	"net"
	"log"
	"io/ioutil"
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"io"
)

const protocol = "tcp"
const nodeVersion = 1
const commandLength = 12

var nodeAddress string
var miningAddress string
var KnownNodes = []string{"localhost:4000"}
var blocksInTransit = [][]byte{}
var mempool = make(map[string]Transaction)

type version struct {
	Version			int
	BestHeight	int
	AddrFrom		string
}

type getblocks struct {
	AddrFrom string
}

type inv struct {
	AddrFrom	string
	Type			string
	Items			[][]byte
}

type getdata struct {
	AddrFrom	string
	Type			string
	ID				[]byte
}

type block struct {
	AddrFrom	string
	Block			[]byte
}

type tx struct {
	AddrFrom		string
	Transaction	[]byte
}

func StartServer(nodeID, minerAddress string) {
	nodeAddress = fmt.Sprintf("localhost:%s", nodeID)
	miningAddress = minerAddress
	ln, err := net.Listen(protocol, nodeAddress)
	if err != nil {
		log.Panic(err)
	}
	defer ln.Close()

	bc := NewBlockchain(nodeID)

	if nodeAddress != KnownNodes[0] {
		sendVersion(KnownNodes[0], bc)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		go handleConnection(conn, bc)
	}
}

func sendVersion(addr string, bc *Blockchain) {
	bestHeight := bc.GetBestHeight()
	payload := gobEncode(version{nodeVersion, bestHeight, nodeAddress})

	request := append(commandToBytes("version"), payload...)

	sendData(addr, request)
}

func sendData(addr string, data []byte) {
	conn, err := net.Dial(protocol, addr)
	if err != nil {
		fmt.Printf("%s is not available\n", addr)
		var updateNodes []string

		for _, node := range KnownNodes {
			if node != addr {
				updateNodes = append(updateNodes, node)
			}
		}

		KnownNodes = updateNodes
		return
	}
	defer conn.Close()

	_, err = io.Copy(conn, bytes.NewReader(data))
	if err != nil {
		log.Panic(err)
	}
}

func gobEncode(data interface{}) []byte {
	var buff bytes.Buffer

	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}

func commandToBytes(command string) []byte {
	var bytes [commandLength]byte

	for i, c := range command {
		bytes[i] = byte(c)
	}

	return bytes[:]
}

func bytesToCommand(bytes []byte) string {
	var command []byte

	for _, b := range bytes {
		if b != 0x0 {
			command = append(command, b)
		}
	}

	return string(command)
}

func handleConnection(conn net.Conn, bc *Blockchain) {
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Panic(err)
	}
	command := bytesToCommand(request[:commandLength])
	fmt.Printf("Received %s command\n", command)

	switch command {
		case "version":
			handleVersion(request, bc)
		case "getblocks":
			handleGetBlocks(request, bc)
		case "inv":
			handleInv(request, bc)
		case "getdata":
			handleGetData(request, bc)
		case "block":
			handleBlock(request, bc)
		case "tx":
			handleTx(request, bc)
		default:
			fmt.Println("Unknown command!")
	}

	conn.Close()
}

func handleVersion(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload version

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	localBestHeight := bc.GetBestHeight()
	foreignBestHeight := payload.BestHeight

	if localBestHeight < foreignBestHeight {
		sendGetBlocks(payload.AddrFrom)
	} else if localBestHeight > foreignBestHeight {
		sendVersion(payload.AddrFrom, bc)
	}

	if !nodeIsKnown(payload.AddrFrom) {
		KnownNodes = append(KnownNodes, payload.AddrFrom)
	}
}

func nodeIsKnown(addr string) bool {
	for _, node := range KnownNodes {
		if node == addr {
			return true
		}
	}

	return false
}

func handleGetBlocks(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload getblocks

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blocks := bc.GetBlockHashes()
	sendInv(payload.AddrFrom, "block", blocks)
}

func sendGetBlocks(address string) {
	payload := gobEncode(getblocks{nodeAddress})
	request := append(commandToBytes("getblocks"), payload...)

	sendData(address, request)
}

func sendInv(address, kind string, items [][]byte) {
	inventory := inv{nodeAddress, kind, items}
	payload := gobEncode(inventory)
	request := append(commandToBytes("inv"), payload...)

	sendData(address, request)
}

func handleInv(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload inv

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("Received inventory with %d %s\n", len(payload.Items), payload.Type)

	if payload.Type == "block" {
		blocksInTransit = payload.Items

		blockHash := payload.Items[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	}

	if payload.Type == "tx" {
		txID := payload.Items[0]

		if mempool[hex.EncodeToString(txID)].ID == nil {
			sendGetData(payload.AddrFrom, "tx", txID)
		}
	}
}

func sendGetData(address, kind string, id []byte) {
	payload := gobEncode(getdata{nodeAddress, kind, id})
	request := append(commandToBytes("getdata"), payload...)

	sendData(address, request)
}

func handleGetData(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload getdata

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	if payload.Type == "block" {
		block, err := bc.GetBlock(payload.ID)
		if err != nil {
			return
		}

		sendBlock(payload.AddrFrom, &block)
	}

	if payload.Type == "tx" {
		txID := hex.EncodeToString(payload.ID)
		tx := mempool[txID]

		SendTx(payload.AddrFrom, &tx)
	}
}

func sendBlock(addr string, b *Block) {
	data := block{nodeAddress, b.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("block"), payload...)

	sendData(addr, request)
}

func SendTx(addr string, transaction *Transaction) {
	data := tx{nodeAddress, transaction.Serialize()}
	payload := gobEncode(data)
	request := append(commandToBytes("tx"), payload...)

	sendData(addr, request)
}

func handleBlock(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload block

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	blockData := payload.Block
	block := DeserializeBlock(blockData)

	fmt.Println("Received a new block!")
	bc.AddBlock(block)

	fmt.Printf("Added block %x\n", block.Hash)

	if len(blocksInTransit) > 0 {
		blockHash := blocksInTransit[0]
		sendGetData(payload.AddrFrom, "block", blockHash)

		blocksInTransit = blocksInTransit[1:]
	} else {
		UTXOSet := UTXOSet{bc}
		UTXOSet.Reindex()
	}
}

func handleTx(request []byte, bc *Blockchain) {
	var buff bytes.Buffer
	var payload tx

	buff.Write(request[commandLength:])
	dec := gob.NewDecoder(&buff)
	err := dec.Decode(&payload)
	if err != nil {
		log.Panic(err)
	}

	txData := payload.Transaction
	tx := DeserializeTransaction(txData)
	mempool[hex.EncodeToString(tx.ID)] = tx

	if nodeAddress == KnownNodes[0] {
		for _, node := range KnownNodes {
			if node != nodeAddress && node != payload.AddrFrom {
				sendInv(node, "tx", [][]byte{tx.ID})
			}
		}
	} else {
		if len(mempool) >= 1 && len(miningAddress) > 0 {
			MineTransactions:
				var txs []*Transaction

				for id := range mempool {
					tx := mempool[id]
					txs = append(txs, &tx)
				}

				if len(txs) == 0 {
					fmt.Println("All transactions are invalid! Waiting for new ones...")
					return
				}

				cbTx := NewCoinbaseTX(miningAddress, "")
				txs = append(txs, cbTx)

				newBlock := bc.MineBlock(txs)
				UTXOSet := UTXOSet{bc}
				UTXOSet.Update(newBlock)

				fmt.Println("New block is mined!")

				for _, tx := range txs {
					txID := hex.EncodeToString(tx.ID)
					delete(mempool, txID)
				}

				for _, node := range KnownNodes {
					if node != nodeAddress {
						sendInv(node, "block", [][]byte{newBlock.Hash})
					}
				}

				if len(mempool) > 0 {
					goto MineTransactions
				}
		}
	}
}
