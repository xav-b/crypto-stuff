package main

import "bytes"

// TXInput represents a transaction input
type TXInput struct {
	// Txid stores the ID of a previous transaction
	Txid []byte
	// Vout is an index of an output within that stransaction
	Vout      int
	Signature []byte
	// Raw public key (not hashed)
	PubKey []byte
}

// UsesKey checks whether the address initiated the transaction
func (in *TXInput) UsesKey(pubKeyHash []byte) bool {
	lockingHash := HashPubKey(in.PubKey)

	return bytes.Compare(lockingHash, pubKeyHash) == 0
}
