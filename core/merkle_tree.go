package core

import (
	"crypto/sha256"
)

type MerkleTree struct {
	Root *MerkleNode
}

type MerkleNode struct {
	Left *MerkleNode
	Right *MerkleNode
	Data []byte
}

func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	node := MerkleNode{}

	if left == nil && right == nil {
		hash := sha256.Sum256(data)
		node.Data = hash[:]
	} else {
		hash := append(left.Data, right.Data...)
		newHash := sha256.Sum256(hash)
		node.Data = newHash[:]
	}

	node.Left = left
	node.Right = right

	return &node
}

func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []*MerkleNode

	if len(data) % 2 == 1 {
		data = append(data, data[len(data) - 1])
	}

	for _, dat := range data {
		node := NewMerkleNode(nil, nil, dat)
		nodes = append(nodes, node)
	}
	
	for i := 0; i < len(data) / 2; i++ {
		var next []*MerkleNode
		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(nodes[j], nodes[j + 1], nil)
			next = append(next, node)
		}
		nodes = next
	}

	tree := MerkleTree{nodes[0]}

	return &tree
}

