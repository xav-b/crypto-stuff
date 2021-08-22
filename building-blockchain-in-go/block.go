package main

import (
	"bytes"
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
	// Data is the actual valuable information containing in the block
	Data          []byte
	PrevBlockHash []byte
	// Hash is the succcessful hash computed by the PoW
	Hash []byte
	// we also save the nonce so it's possible to verify the PoW
	Nonce int
}

func MineBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{BLOCK_VERSION, time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}

	pow := NewProofOfWork(block)
	nonce, hash := pow.Mine()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
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
