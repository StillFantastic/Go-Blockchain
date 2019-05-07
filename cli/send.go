package cli

import (
	"fmt"
	"github.com/StillFantastic/go-blockchain/core"
	"github.com/StillFantastic/go-blockchain/wallet"
	"log"
)

func (cli *CLI) send(from, to string, amount int, nodeID string, mine bool) {
	if !wallet.ValidateAddress(from) {
		log.Panic("ERROR: Sender address is not valid")
	}
	if !wallet.ValidateAddress(to) {
		log.Panic("ERROR: Recipient address is not valid")
	}

	bc := core.NewBlockchain(nodeID)
	UTXOSet := core.UTXOSet{bc}
	defer bc.DB.Close()

	wallets, err := wallet.NewWallets(nodeID)
	if err != nil {
		log.Panic(err)
	}
	wt := wallets.GetWallet(from)
	tx := bc.NewUTXOTransaction(&wt, to, amount, &UTXOSet)

	if bc.VerifyTransaction(tx) == false {
		fmt.Println("#########")
	}

	if mine {
		cbTx := core.NewCoinbaseTX(from, "")
		txs := []*core.Transaction{cbTx, tx}
		newBlock := bc.MineBlock(txs)
		UTXOSet.Update(newBlock)
	} else {
		core.SendTx(core.KnownNodes[0], tx)
	}

	fmt.Println("Success!")
}
