package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"log"
	"time"
)

const BLOCK_VERSION = 1

// Block is a simplified implementation of what is described in Bitcoin
type Block struct {
	// Block version number
	Version int
	// Timestamp is when the block is created
	Timestamp int64

	Transactions []*Transaction

	PrevBlockHash []byte
	// Hash is the succcessful hash computed by the PoW
	Hash []byte
	// we also save the nonce so it's possible to verify the PoW
	Nonce int
}

func MineBlock(transactions []*Transaction, prevBlockHash []byte) *Block {
	block := &Block{BLOCK_VERSION, time.Now().Unix(), transactions, prevBlockHash, []byte{}, 0}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Mine()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// NewGenesisBlock creates and returns genesis Block
func MineGenesisBlock(coinbase *Transaction) *Block {
	return MineBlock([]*Transaction{coinbase}, []byte{})
}

// Serialize translates all block information into a format easy to store or
// transfer
func (b *Block) Serialize() []byte {
	var result bytes.Buffer
	// we pick the `gob` library as it is part of the standard library and does
	// the work good enough for our pet implementation. Valid alternatives could
	// be JSON, protocol buffers, ...
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(b)
	if err != nil {
		log.Panic(err)
	}

	return result.Bytes()
}

func DeserializeBlock(d []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block
}

// TODO: implement:
// Bitcoin uses a more elaborate technique: it represents all transactions
// containing in a block as a Merkle tree and uses the root hash of the tree in
// the Proof-of-Work system. This approach allows to quickly check if a block
// contains certain transaction, having only just the root hash and without
// downloading all the transactions.
func (b *Block) HashTransactions() []byte {
	// hashes of each transactions
	var txHashes [][]byte
	// final 32 bits hash of the hashes
	var txHash [32]byte

	for _, tx := range b.Transactions {
		txHashes = append(txHashes, tx.Hash())
	}
	txHash = sha256.Sum256(bytes.Join(txHashes, []byte{}))

	return txHash[:]
}
