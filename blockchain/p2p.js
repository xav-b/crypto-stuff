const debug = require('debug')('blockchain:p2p')
const wrtc = require('wrtc')
const Exchange = require('peer-exchange')
const net = require('net')
const uuidv4 = require('uuid/v4')

// NOTE so far we didn't need to use blockchain or block
const Block = require('./block')
const Messages = require('./messages')
const {
  REQUEST_LATEST_BLOCK,
  RECEIVE_LATEST_BLOCK,
  REQUEST_BLOCKCHAIN,
  RECEIVE_BLOCKCHAIN,
  REQUEST_PEERS,
} = require('./messages/types')

const p2p = new Exchange('blockchain.js', { wrtc })

class PeerToPeer {
  constructor(blockchain) {
    this.peers = []
    this.blockchain = blockchain
    this.nodeID = uuidv4().replace('-', '')
  }

  serve(port = 3000) {
    const server = net
      .createServer(socket =>
        p2p.accept(socket, (err, conn) => {
          if (err) {
            debug(`❗  ${err}`)
          } else {
            debug('new peer joined the server')
            this.initConnection.call(this, conn)
          }
        })
      )
      .listen(port)
    debug(`listening to peers on ${server.address().address}:${server.address().port}... `)
  }

  // TODO discoverPeers() ?

  join(host, port, cb = null) {
    const socket = net.connect(port, host, () =>
      p2p.connect(socket, (err, conn) => {
        if (err) {
          debug(`❗  ${err}`)
        } else {
          debug('Successfully connected to a new peer!')
          const toCall = cb ? cb : this.initConnection
          toCall.call(this, conn)
        }
      })
    )
  }

  initConnection(conn) {
    this.peers.push(conn)
    this.initMessageHandler(conn)
    this.initErrorHandler(conn)

    // trigger the sync dance
    this.write(conn, Messages.getLatestBlock())
  }

  broadcast(message) {
    this.peers.forEach(peer => this.write(peer, message))
  }

  write(peer, message) {
    message.node = this.nodeID
    peer.write(JSON.stringify(message))
  }

  initErrorHandler(connection) {
    connection.on('error', error => debug(`❗  ${error}`))
  }

  initMessageHandler(connection) {
    connection.on('data', data => {
      const message = JSON.parse(data.toString('utf8'))
      this.handleMessage(connection, message)
    })
  }

  handleMessage(peer, message) {
    switch (message.type) {
      case REQUEST_LATEST_BLOCK:
        debug('peer requested for latest block')
        this.write(peer, Messages.sendLatestBlock(this.blockchain.latestBlock))
        break
      case REQUEST_BLOCKCHAIN:
        debug('peer requested for blockchain')
        this.write(peer, Messages.sendBlockchain(this.blockchain.get()))
        break
      case RECEIVE_LATEST_BLOCK:
        debug('received latest block')
        this.handleReceivedLatestBlock(message, peer)
        break
      case RECEIVE_BLOCKCHAIN:
        debug('received blockchain')
        this.handleReceivedBlockchain(message)
        break
      case REQUEST_PEERS:
        debug('peer requested network details')
        this.write(peer, Messages.sendPeers(p2p))
        break
      default:
        debug(`❓  Received unknown message type ${message.type}`)
    }
  }

  handleReceivedLatestBlock(message, peer) {
    // NOTE I don't what kind of stuff I can get from `peer`
    const receivedBlock = message.data

    debug(`Peer sent over blockchain [${receivedBlock.index}:${receivedBlock.hash}]`)
    return this.syncBlockchain([receivedBlock])
  }

  handleReceivedBlockchain(message) {
    const receivedBlocks = message.data.sort((b1, b2) => b1.index - b2.index)

    debug(`Peer sent over blockchain (${receivedBlocks.length})`)
    return this.syncBlockchain(receivedBlocks)
  }

  syncBlockchain(receivedBlocks) {
    const latestBlockReceived = receivedBlocks[receivedBlocks.length - 1]
    const latestBlockHeld = this.blockchain.latestBlock

    if (latestBlockReceived.index <= latestBlockHeld.index) {
      debug('received latest block is not longer than current blockchain. Do nothing')
      // NOTE do something with return value ?
      return null
    }

    debug(
      `blockchain possibly behind. Received latest block is #${latestBlockReceived.index}. Current latest block is #${latestBlockHeld.index}.`
    )

    if (latestBlockHeld.hash === latestBlockReceived.previousHash) {
      debug('Previous hash received is equal to current hash. Append received block to blockchain')
      // FIXME addBlockFromPeer not defined !
      this.blockchain.addBlock(
        new Block(
          latestBlockReceived.index,
          latestBlockReceived.previousHash,
          latestBlockReceived.hash,
          latestBlockReceived.data,
          latestBlockReceived.nonce,
          latestBlockReceived.timestamp
        )
      )
      this.broadcast(Messages.sendLatestBlock(this.blockchain.latestBlock))
    } else if (receivedBlocks.length === 1) {
      debug('received previous hash different from current hash. Get entire blockchain from peer')
      this.broadcast(Messages.getBlockchain())
    } else {
      debug('peer blockchain is longer than current blockchain')
      this.blockchain.replaceChain(receivedBlocks)
      this.broadcast(Messages.sendLatestBlock(this.blockchain.latestBlock))
    }
  }

  closeConnection() {
    p2p.close(err => debug(`❗  ${err}`))
  }
}

module.exports = PeerToPeer
