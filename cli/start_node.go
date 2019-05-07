package cli

import (
	"fmt"
	"github.com/StillFantastic/go-blockchain/wallet"
	"github.com/StillFantastic/go-blockchain/core"
	"log"
)

func (cli *CLI) startNode(nodeID, minerAddress string) {
	fmt.Printf("Starting node %s\n", nodeID)
	if len(minerAddress) > 0 {
		if wallet.ValidateAddress(minerAddress) {
			fmt.Println("Mining is on. Address to receive rewards: ", minerAddress)
		} else {
			log.Panic("Wrong miner address")
		}
	}
	core.StartServer(nodeID, minerAddress)
}




