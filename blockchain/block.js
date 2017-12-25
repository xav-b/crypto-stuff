const debug = require('debug')('blockchain:block')

class Block {
  constructor(
    index,
    previousHash,
    hash,
    data = '',
    nonce = 0,
    timestamp = new Date().getTime() / 1000
  ) {
    debug('new Block created [%d:%s]', index, hash)

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
    this.data = data
    // in bitcoin this.transactions = data

    // How many iterations did we go through before we found a valid block
    this.nonce = nonce
  }

  static get genesis() {
    return new Block(
      0,
      '0',
      // TODO generate it
      '000dc75a315c77a1f9c98fb6247d03dd18ac52632d7dc6a9920261d8109b37cf',
      '-- genesis block --',
      // random nonce
      Math.floor(Math.random() * 10000)
    )
  }
}

module.exports = Block
