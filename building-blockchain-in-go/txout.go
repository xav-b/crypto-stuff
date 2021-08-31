package main

import "bytes"

// TXOutput represents a transaction output
type TXOutput struct {
	// Value is the actual storage of coins
	// in satoshis (== 0.00000001 BTC)
	Value int
	// PubKeyHash is the hash of the public key that can unlock the output
	// Script reference: https://en.bitcoin.it/wiki/Script
	PubKeyHash []byte
}

// Lock signs the output
func (out *TXOutput) Lock(address []byte) {
	pubKeyHash := Base58Decode(address)
	pubKeyHash = pubKeyHash[1 : len(pubKeyHash)-4]
	out.PubKeyHash = pubKeyHash
}

// IsLockedWithKey checks if the output can be used by the owner of the pubkey
// This is very similar to TXInput.UsesKey but the public key used is already
// hashed, while TXInput stores the raw public key
func (out *TXOutput) IsLockedWithKey(pubKeyHash []byte) bool {
	return bytes.Compare(out.PubKeyHash, pubKeyHash) == 0
}

// NewTXOutput create a new TXOutput
func NewTXOutput(value int, address string) *TXOutput {
	txo := &TXOutput{value, nil}
	txo.Lock([]byte(address))

	return txo
}
