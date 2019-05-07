package cli

import (
	"fmt"
	"github.com/StillFantastic/go-blockchain/core"
	"github.com/StillFantastic/go-blockchain/tool"
)

func (cli *CLI) getBalance(address, nodeID string) {
	bc := core.NewBlockchain(nodeID)
	UTXOSet := core.UTXOSet{bc}
	defer bc.DB.Close()

	balance := 0

	pubKeyHash := tool.Base58Decode([]byte(address))
	pubKeyHash = pubKeyHash[1:len(pubKeyHash) - 4]
	UTXOs := UTXOSet.FindUTXO(pubKeyHash)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of %s = %d\n", address, balance)
}
