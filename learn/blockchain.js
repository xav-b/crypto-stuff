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

class Blockchain {
  constructor (difficulty = 3) {
    this.blockchain = [Block.genesis()]
    this.difficulty = difficulty
  }

  get() {
    return this.blockchain
  }

  get latestBlock() {
    return this.blockchain[this.blockchain.length - 1]
  }

  isValidHash(hash) {
    for (var i = 0; i < hash.length; i++) {
      if (hash[i] !== '0') {
        break
      }
    }
    return i >= this.difficulty
  }

  calculateHashForBlock(block) {
    const { index, previousHash, timestamp, transactions, nonce } = block
    return this.calculateHash(
      index,
      previousHash,
      timestamp,
      transactions,
      nonce
    )
  }

  calculateHash(index, previousHash, timestamp, data, nonce) {
    return crypto
      .createHash('sha256')
      .update(index + previousHash + timestamp + data + nonce)
      .digest('hex')
  }

  generateNextBlock(data) {
    const nextIndex = this.latestBlock.index + 1
    const previousHash = this.latestBlock.hash
    let timestamp = new Date().getTime()
    let nonce = 0
    let nextHash = this.calculateHash(nextIndex, previousHash, timestamp, data, nonce)

    while (!this.isValidHashDifficulty(nextHash)) {
      nonce = nonce + 1
      timestamp = new Date().getTime()
      nextHash = this.calculateHash(nextIndex, previousHash, timestamp, data, nonce)
    }

    const nextBlock = new Block(
      nextIndex,
      previousBlock.hash,
      nextTimestamp,
      data,
      nextHash,
      nonce
    )

    return nextBlock
  }

  addBlock(newBlock) {
    if (this.isValidNewBlock(newBlock, this.latestBlock)) {
      this.blockchain.push(newBlock)
    } else {
      throw 'Error: Invalid block'
    }
  }

  isValidNextBlock(nextBlock, previousBlock) {
    const nextBlockHash = this.calculateHashForBlock(nextBlock);

    if (previousBlock.index + 1 !== nextBlock.index) {
      return false
    } else if (previousBlock.hash !== nextBlock.previousHash) {
      return false
    } else if (nextBlockHash !== nextBlock.hash) {
      return false
    } else if (!this.isValidHashDifficulty(nextBlockHash)) {
      return false
    } else {
      return true
    }
  }

  mine(data) {
    const newBlock = this.generateNextBlock(data);
    try {
      this.addBlock(newBlock)
    } catch (err) {
      throw err;
    }
  }

  isValidChain(chain) {
    if (JSON.stringify(chain[0]) !== JSON.stringify(Block.genesis)) {
      return false
    }

    const tempChain = [chain[0]];
    for (let i = 1; i < chain.length; i = i + 1) {
      if (this.isValidNextBlock(chain[i], tempChain[i - 1])) {
        tempChain.push(chain[i])
      } else {
        return false
      }
    }

    return true
  }

  isChainLonger(chain) {
    return chain.length > this.blockchain.length
  }

  replaceChain(newChain) {
    if (this.isValidChain(newChain) && this.isChainLonger(newChain)) {
      this.blockchain = JSON.parse(JSON.stringify(newChain))
    } else {
      throw 'Error: invalid chain'
    }
  }
}

export default { Block, Blockchain }
