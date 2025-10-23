/**
 * Credits:
 * https://medium.com/@lhartikk/a-blockchain-in-200-lines-of-code-963cc1cc0e54
 */

const debug = require('debug')('blockchain:core')
const CryptoJS = require('crypto-js')

const Block = require('./block')

// NOTE what with all those toString ?
// NOTE method of Block ?
const calculateHash = (index, previousHash, timestamp, data, nonce) =>
  CryptoJS.SHA256(index + previousHash + timestamp + data + nonce).toString()

function isValidHash(hash, difficulty) {
  let i
  for (i = 0; i < hash.length; i++) {
    if (hash[i] !== '0') {
      break
    }
  }
  return i >= difficulty
}

function isValidNextBlock(nextBlock, previousBlock, difficulty) {
  const { index, previousHash, timestamp, data, nonce } = nextBlock
  const nextBlockHash = calculateHash(index, previousHash, timestamp, data, nonce)

  if (previousBlock.index + 1 !== nextBlock.index) {
    debug('❌  new block has invalid index')
    return false
  } else if (previousBlock.hash !== nextBlock.previousHash) {
    debug('❌  new block has invalid previous hash')
    return false
  } else if (nextBlockHash !== nextBlock.hash) {
    debug(`❌  invalid hash: ${nextBlockHash} ${nextBlock.hash}`)
    return false
  } else if (!isValidHash(nextBlockHash, difficulty)) {
    debug(`❌  hash does not meet difficulty requirements: ${nextBlockHash}`)
    return false
  }

  return true
}

class Blockchain {
  constructor(difficulty = 3) {
    debug('initializing a new blockchain')
    this.blockchain = [Block.genesis]
    this.difficulty = difficulty
  }

  get() {
    return this.blockchain
  }

  get latestBlock() {
    return this.blockchain[this.blockchain.length - 1]
  }

  mine(data) {
    const newBlock = this.generateNextBlock(data)
    if (this.addBlock(newBlock)) {
      debug('a new block was successfully mined', newBlock)
    }
  }

  replaceChain(newBlocks) {
    if (!this.isValidChain(newBlocks)) {
      debug("❌  replacement chain is not valid. Won't replace existing blockchain.")
      return null
    }

    if (newBlocks.length <= this.blockchain.length) {
      debug("❌  Replacement chain is shorter than original. Won't replace existing blockchain.")
      return null
    }

    debug('✅  Received blockchain is valid. Replacing current blockchain with received blockchain')
    this.blockchain = newBlocks.map(
      json =>
        new Block(json.index, json.previousHash, json.hash, json.data, json.nonce, json.timestamp)
    )
  }

  generateNextBlock(blockData) {
    const previousBlock = this.latestBlock
    const nextIndex = previousBlock.index + 1
    const nextTimestamp = new Date().getTime() / 1000
    let nonce = 0
    let nextHash = ''

    debug('mining a new block - index=%d', nextIndex)
    while (!isValidHash(nextHash, this.difficulty)) {
      nonce = nonce + 1
      nextHash = calculateHash(nextIndex, previousBlock.hash, nextTimestamp, blockData, nonce)
      // TODO nice UI which replace the computed hash, with a spinner
    }

    return new Block(nextIndex, previousBlock.hash, nextHash, blockData, nonce, nextTimestamp)
  }

  addBlock(newBlock) {
    if (isValidNextBlock(newBlock, this.latestBlock, this.difficulty)) {
      this.blockchain.push(newBlock)
      return true
    }

    return false
  }

  isValidChain(chain) {
    if (JSON.stringify(chain[0]) !== JSON.stringify(Block.genesis)) {
      return false
    }

    const tempChain = [chain[0]]
    for (let i = 1; i < chain.length; i = i + 1) {
      if (isValidNextBlock(chain[i], tempChain[i - 1], this.difficulty)) {
        tempChain.push(chain[i])
      } else {
        return false
      }
    }

    return true
  }
}

module.exports = new Blockchain()
