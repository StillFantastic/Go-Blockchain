# Blockchain implementation in Golang

## Succesfully implemented featurre
* Public-key cryptography
  >Used to generate bitcoin address and digital signature (secp256k1's elliptic curve)
* Proof of work
	>The consensus algorithm
* UTXO (Unspent Transaction Output)
	>The way how transactions work
* Merkle Tree
  >A data structure that stores transactions


## P2P Network
Considering the complexity of the implementation, I used an alternative.
Dividing nodes into three catagories:
1. Wallet node 
2. Central node
3. Miner node  

Wallet node is responsible for generating transacions, miner node is responsible for mining, central node is responsible for transfering the data between Wallet node and miner node.

Here we can open 3 terminal window to simulate the operation of each node.

1. Enter the following command in each window
	```
	export NODE_ID=4000   (Central)
	export NODE_ID=4001   (Wallet)
	export NODE_ID=4002   (Miner)
	```
    This allows them to use different ports to simulate the communication between nodes.

2. In central node terminal:
	```
	./btc createwallet
	./btc createblockchain -address CENTRAL_ADDR
	cp blockchain_4000.db blockchain_4001.db
	cp blockchain_4000.db blockchain_4002.db
	```
    Generate new address and genesis block, and send the mining reward to the new address and send the block to miner node and wallet node.  

	  (In real bitcoin system, the genesis block is hardcoded.)

3. In wallet terminal:
	```
	./btc createwallet
	./btc createwallet
	```
    Generate two new addresses.

4. In central terminal:
	```
	./btc send -from CENTRAL_ADDR -to WALLET1_ADDR -amount 10 -mine
	./btc startnode
	```
	Make a transaction, and then start running the node, don't close it.

5. In wallet terminal:
	```
	./btc startnode
	```
    This command syncs wallet node and central node, and then close it.
	```
	./btc getbalance -address WALLET1_ADDR
	```
    Make sure it really receives 10 coins.
	
6. In miner terminal:
	```
	./btc createwallet
	./btc startnode -miner MINER_ADDR
	```
    Start miner node.

7. In wallet terminal
	```
	./btc send -from WALLET1_ADDR -to WALLET2_ADDR -amount 10
	```

8. In the miner terminal, we can see it receives a new transaction and is mining a new block.

9. In wallet terminal
	```
	./btc startnode
	```
    Start the node to receive the latest block.

10. In wallet nodde
	```
	./btc getbalance -address WALLET2_ADDR
	```
    Now it has 10 coins!


