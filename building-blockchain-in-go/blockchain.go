package main

import (
	"log"

	"github.com/boltdb/bolt"
)

const (
	DB_FILE = "blockchain.db"
	// bitcoin (for exmaple) stores 4 different entites but at this stage blocks
	// are the only bits of data to be persisted
	BLOCKS_BUCKET = "blocks"
)

// Blockchain Iterator lets us go through the saved blockchain, in a way wich is
// ordered (by the chain of blocks) and efficient (without loading all blocks in
// memory)
type BlockchainIterator struct {
	// currentHash is the pointer to the current block in the iteration
	currentHash []byte
	db          *bolt.DB
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	// start at the tip and walk toward the oldest block
	//
	// Note that a valid blockchain is defined as the longest one. Therefore
	// picking the tip is like `voting` for what we considere to be the valid
	// blockchain, and not some (hopefully temporary) forks
	return &BlockchainIterator{bc.tip, bc.db}
}

// Next yields the next block in the blockchain
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	_ = i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKS_BUCKET))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	// point at the next (older) block in the chain
	i.currentHash = block.PrevBlockHash

	return block
}

type Blockchain struct {
	// hash of the latest block
	tip []byte
	// blocks DB
	db *bolt.DB
}

func (bc *Blockchain) AddBlock(data string) {
	var lastHash []byte

	// open a read-only transaction
	_ = bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKS_BUCKET))
		// get latest block hash
		lastHash = b.Get([]byte("l"))

		return nil
	})

	newBlock := MineBlock(data, lastHash)

	// save the new block
	_ = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKS_BUCKET))
		_ = b.Put(newBlock.Hash, newBlock.Serialize())
		_ = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash

		return nil
	})
}

// NewGenesisBlock creates the first block of the chain
func NewGenesisBlock() *Block {
	return MineBlock("Genesis Block", []byte{})
}

func NewBlockchain() *Blockchain {
	// tip of the blockchain
	var tip []byte

	log.Printf("opening blockcain db: %s\n", DB_FILE)
	db, err := bolt.Open(DB_FILE, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	// start a read/write boltdb transaction
	err = db.Update(func(tx *bolt.Tx) error {
		// load the blocks bucket within the blockchain database
		b := tx.Bucket([]byte(BLOCKS_BUCKET))

		if b == nil {
			// no blocks saved in this blockchain db
			// let's initialise a new blockchain, and therefore mine the Genesis block
			genesis := NewGenesisBlock()

			// initialise the DB and store our first block
			b, _ := tx.CreateBucket([]byte(BLOCKS_BUCKET))
			// store the serialized block, indexed at his hash
			_ = b.Put(genesis.Hash, genesis.Serialize())
			// store the tip of the blockchain
			_ = b.Put([]byte("l"), genesis.Hash)
			tip = genesis.Hash
		} else {
			// found an existing blockchain, set the tip of it
			tip = b.Get([]byte("l"))
		}

		return nil
	})

	bc := Blockchain{tip, db}

	return &bc
}
