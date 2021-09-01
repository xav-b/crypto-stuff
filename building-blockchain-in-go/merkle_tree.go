package main

import (
	"crypto/sha256"
	"log"
)

// MerkleTree represent a Merkle tree
type MerkleTree struct {
	RootNode *MerkleNode
}

// MerkleNode represent a Merkle tree node
type MerkleNode struct {
	Left  *MerkleNode
	Right *MerkleNode
	Data  []byte
}

// NewMerkleTree creates a new Merkle tree from a sequence of data
func NewMerkleTree(data [][]byte) *MerkleTree {
	var nodes []MerkleNode

	if len(data)%2 != 0 {
		// odd number of leaves
		log.Println("copying last leaf node to have an even number")
		data = append(data, data[len(data)-1])
	}

	log.Println("initialising tree bottom level")
	for _, datum := range data {
		node := NewMerkleNode(nil, nil, datum)
		nodes = append(nodes, *node)
	}

	log.Printf("assembling the remaining %d levels\n", len(data)/2)
	for i := 0; i < len(data)/2; i++ {
		var newLevel []MerkleNode

		// for each pair of nodes, create one new at this level
		for j := 0; j < len(nodes); j += 2 {
			node := NewMerkleNode(&nodes[j], &nodes[j+1], nil)
			newLevel = append(newLevel, *node)
		}

		// overwrite with the new level
		nodes = newLevel
	}

	// there should be only one leaf node, at the top of the tree
	return &MerkleTree{&nodes[0]}
}

// NewMerkleNode creates a new Merkle tree node
func NewMerkleNode(left, right *MerkleNode, data []byte) *MerkleNode {
	mNode := MerkleNode{}

	// compute node data
	if left == nil && right == nil {
		// first leaf of the tree, initialise from external data
		hash := sha256.Sum256(data)
		mNode.Data = hash[:]
	} else {
		// a new node ot of 2 leaves is the hash of their 2 respective hashes
		prevHashes := append(left.Data, right.Data...)
		hash := sha256.Sum256(prevHashes)
		mNode.Data = hash[:]
	}

	// link to right and left nodes (if any)
	mNode.Left = left
	mNode.Right = right

	return &mNode
}
