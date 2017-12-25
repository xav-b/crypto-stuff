const debug = require('debug')('blockchain:messages')
const {
  REQUEST_LATEST_BLOCK,
  RECEIVE_LATEST_BLOCK,
  REQUEST_BLOCKCHAIN,
  RECEIVE_BLOCKCHAIN,
  REQUEST_PEERS,
  RECEIVE_PEERS,
} = require('./types')

// TODO seperate module, those static serves nothing
class Messages {
  static getLatestBlock() {
    return { type: REQUEST_LATEST_BLOCK }
  }

  // NOTE could be sendBlock
  static sendLatestBlock(block) {
    debug('sending latest block [%d:%s]', block.index, block.hash)
    return {
      type: RECEIVE_LATEST_BLOCK,
      data: block,
    }
  }

  static getBlockchain() {
    return { type: REQUEST_BLOCKCHAIN }
  }

  static sendBlockchain(blockchain) {
    return { type: RECEIVE_BLOCKCHAIN, data: blockchain }
  }

  static getPeers() {
    return { type: REQUEST_PEERS }
  }

  static sendPeers(exchange) {
    return { type: RECEIVE_PEERS, data: { nodes: exchange.peers.length } }
  }
}

module.exports = Messages
