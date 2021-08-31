package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
)

const (
	// Mining difficulty
	// We will need to find a hash value lower than a set target. This target is
	// computed as a 256bits hash, with the first `24 / 8` bits set to 0.
	//
	// So we can increase the difficulty by asking for more leading zeros, i.e.
	// by increasing the `targetBits` value by steps of 8. And vice-e-versa:
	// targetBits=16 will only need the PoW to figure out a hash with 2 leading
	// zeros.
	//
	// In Bitcoin, "target bits" is the block header storing the difficulty at which
	// the block was mined. Unlike bitcoin though, this is not dynamically adjusted
	// to miners capacity
	// targetBits = 24
	targetBits = 16 // FIXME: it's too easy, only for dev
	// set a large upper boundary to our infinite loop
	MAX_NONCE = math.MaxInt64
)

type ProofOfWork struct {
	block *Block

	// our proof of work consists of finding a hash from block's data +
	// something, which is lower than the target
	target *big.Int
}

func NewProofOfWork(b *Block) *ProofOfWork {
	// initialise to 1 and shift it left by `256 - targetBits` bits
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	pow := &ProofOfWork{b, target}

	return pow
}

func (pow *ProofOfWork) prepareData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			// block data
			pow.block.PrevBlockHash,
			pow.block.HashTransactions(),
			IntToHex(pow.block.Timestamp),
			// pow properties
			IntToHex(int64(targetBits)),
			// nonce here is the counter from the Hashcash algo
			IntToHex(int64(nonce)),
		},
		[]byte{},
	)

	return data
}

func (pow *ProofOfWork) Mine() (int, []byte) {
	var hashInt big.Int
	var hash [32]byte
	nonce := 0

	fmt.Printf("Mining the block containing \"%+v\"\n", pow.block.Transactions)
	for nonce < MAX_NONCE {
		// create a byte representation of block's data, nonce and POW target
		data := pow.prepareData(nonce)
		hash = sha256.Sum256(data)
		fmt.Printf("\r%x (nonce: %d)", hash, nonce)
		// convert hash to bigint
		hashInt.SetBytes(hash[:])

		// validate POW
		if hashInt.Cmp(pow.target) == -1 {
			// valid!
			break
		} else {
			// not yet, increment and try again
			nonce++
		}
	}
	fmt.Print("\n\n")

	// return nonce and hash winners
	return nonce, hash[:]
}

// Validate takes a newly minted block and check that its nonce and hash pass
// the PoW test
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.block.Nonce)
	hash := sha256.Sum256(data)
	hashInt.SetBytes(hash[:])

	isValid := hashInt.Cmp(pow.target) == -1

	return isValid
}
