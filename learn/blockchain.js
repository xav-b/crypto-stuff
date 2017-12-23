/**
 * Credits:
 * https://medium.com/@lhartikk/a-blockchain-in-200-lines-of-code-963cc1cc0e54
 */

const CryptoJS = require('crypto-js')

class Block {
  constructor(index, previousHash, timestamp, data, hash, nonce) {
    // Which block is it
    this.index = index
    // When was the block added
    this.timestamp = timestamp
    // Is the block valid
    this.hash = hash.toString()
    // is the previous block valid
    this.previousHash = previousHash.toString()

    // What information is stored on the block. Instead of having text as data,
    // cryptocurrencies have transactions as data
    // Transactions are a record of payment between two parties. When there is
    // an exchange of value, a transaction is created to record it
    this.transactions = data

    // How many iterations did we go through before we found a valid block
    this.nonce = nonce
  }

  get genesis() {
    new Block(
      0,
      '0',
      new Date().getTime(),
      'Welcome to Blockchain Demo 2.0!',
      // TODO generate it
      '000dc75a315c77a1f9c98fb6247d03dd18ac52632d7dc6a9920261d8109b37cf',
      // FIXME what ? random ?
      604
    )
  }
}

// NOTE what with all those toString ?
// NOTE method of Block ?
const calculateHash = (index, previousHash, timestamp, data) =>
  CryptoJS.sha256(index + previousHash + timestamp + data).toString()

export default { Block }
