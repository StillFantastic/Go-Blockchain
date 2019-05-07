package cli

import (
	"fmt"
	"github.com/StillFantastic/go-blockchain/core"
)

func (cli *CLI) createBlockchain(address, nodeID string) {
	bc := core.CreateBlockchain(address, nodeID)
	defer bc.DB.Close()

	UTXOSet := core.UTXOSet{bc}
	UTXOSet.Reindex()
	fmt.Println("Done")
}
