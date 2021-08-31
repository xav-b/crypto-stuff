package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"log"

	"github.com/boltdb/bolt"
)

const (
	DB_FILE = "blockchain.db"
	// bitcoin (for exmaple) stores 4 different entites but at this stage blocks
	// are the only bits of data to be persisted
	BLOCKS_BUCKET = "blocks"

	// actual first Bitcoin message wthin the first transaction
	// check: https://www.blockchain.com/btc/tx/4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b?show_adv=true
	// tutorial: https://medium.com/geekculture/decoding-bitcoins-first-block-coinbase-transaction-aeefe87ceec0
	genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"
)

// Blockchain Iterator lets us go through the saved blockchain, in a way wich is
// ordered (by the chain of blocks) and efficient (without loading all blocks in
// memory)
type BlockchainIterator struct {
	// currentHash is the pointer to the current block in the iteration
	currentHash []byte
	db          *bolt.DB
}

// Next yields the next block in the blockchain
func (i *BlockchainIterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKS_BUCKET))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})

	if err != nil {
		log.Panic(err)
	}

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

func (bc *Blockchain) Iterator() *BlockchainIterator {
	// start at the tip and walk toward the oldest block
	//
	// Note that a valid blockchain is defined as the longest one. Therefore
	// picking the tip is like `voting` for what we considere to be the valid
	// blockchain, and not some (hopefully temporary) forks
	return &BlockchainIterator{bc.tip, bc.db}
}

func (bc *Blockchain) AddBlock(transactions []*Transaction) {
	var lastHash []byte

	for _, tx := range transactions {
		if !bc.VerifyTransaction(tx) {
			log.Panic("ERROR: Invalid transaction")
		}
	}

	// open a read-only transaction
	err := bc.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKS_BUCKET))
		// get latest block hash
		lastHash = b.Get([]byte("l"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := MineBlock(transactions, lastHash)

	// save the new block
	_ = bc.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(BLOCKS_BUCKET))
		_ = b.Put(newBlock.Hash, newBlock.Serialize())
		_ = b.Put([]byte("l"), newBlock.Hash)
		bc.tip = newBlock.Hash

		return nil
	})
}

// NewBlochain loads or initialises a blockchain.
// The address given will receive the award of the geneis block
func NewBlockchain(address string) *Blockchain {
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
			cbtx := NewCoinbaseTX(address, genesisCoinbaseData)
			genesis := MineGenesisBlock(cbtx)

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

func (bc *Blockchain) FindUnspentTransactions(pubKeyHash []byte) []Transaction {
	var unspentTXs []Transaction
	spentTXOs := make(map[string][]int)
	bci := bc.Iterator()

	for {
		// walk down the blockchain
		block := bci.Next()

		for _, tx := range block.Transactions {
			// inspect each transaction
			txID := hex.EncodeToString(tx.ID)

		Outputs:
			// inspect each transactions output
			for outIdx, out := range tx.Vout {
				// Was the output spent?
				if spentTXOs[txID] != nil {
					for _, spentOutIdx := range spentTXOs[txID] {
						if spentOutIdx == outIdx {
							continue Outputs
						}
					}
				}

				// check if this UTXO belongs to the given address
				if out.IsLockedWithKey(pubKeyHash) {
					unspentTXs = append(unspentTXs, *tx)
				}
			}

			if !tx.IsCoinbase() {
				// register the addres's inputs, which by definition have a spent output
				for _, in := range tx.Vin {
					if in.UsesKey(pubKeyHash) {
						inTxID := hex.EncodeToString(in.Txid)
						spentTXOs[inTxID] = append(spentTXOs[inTxID], in.Vout)
					}
				}
			}
		}

		if len(block.PrevBlockHash) == 0 {
			// we reached the genesis block
			break
		}
	}

	return unspentTXs
}

func (bc *Blockchain) FindUTXO(pubKeyHash []byte) []TXOutput {
	var UTXOs []TXOutput
	unspentTransactions := bc.FindUnspentTransactions(pubKeyHash)

	for _, tx := range unspentTransactions {
		for _, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) {
				UTXOs = append(UTXOs, out)
			}
		}
	}

	return UTXOs
}

// FindSpendableOutputs finds and returns unspent outputs to reference in inputs
func (bc *Blockchain) FindSpendableOutputs(pubKeyHash []byte, amount int) (int, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	accumulated := 0

	unspentTXs := bc.FindUnspentTransactions(pubKeyHash)

Work:
	for _, tx := range unspentTXs {
		txID := hex.EncodeToString(tx.ID)

		for outIdx, out := range tx.Vout {
			if out.IsLockedWithKey(pubKeyHash) && accumulated < amount {
				accumulated += out.Value
				unspentOutputs[txID] = append(unspentOutputs[txID], outIdx)

				if accumulated >= amount {
					break Work
				}
			}
		}
	}

	return accumulated, unspentOutputs
}

// FindTransaction finds a transaction by its ID
func (bc *Blockchain) FindTransaction(ID []byte) (Transaction, error) {
	bci := bc.Iterator()

	for {
		block := bci.Next()

		for _, tx := range block.Transactions {
			if bytes.Compare(tx.ID, ID) == 0 {
				return *tx, nil
			}
		}

		if len(block.PrevBlockHash) == 0 {
			break
		}
	}

	return Transaction{}, errors.New("Transaction is not found")
}

func (bc *Blockchain) SignTransaction(tx *Transaction, privKey ecdsa.PrivateKey) {
	prevTXs := make(map[string]Transaction)

	// read all the input's transaction id and fetch the corresponding
	// transaction
	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	tx.Sign(privKey, prevTXs)
}

func (bc *Blockchain) VerifyTransaction(tx *Transaction) bool {
	prevTXs := make(map[string]Transaction)

	for _, vin := range tx.Vin {
		prevTX, err := bc.FindTransaction(vin.Txid)
		if err != nil {
			log.Panic(err)
		}
		prevTXs[hex.EncodeToString(prevTX.ID)] = prevTX
	}

	return tx.Verify(prevTXs)
}
