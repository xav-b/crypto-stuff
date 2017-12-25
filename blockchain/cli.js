#! /usr/bin/env node

const P2P = require('./p2p')
const Messages = require('./messages')
const BlockChain = require('./blockchain')

const command = process.argv[2]

// TODO share messages
function sendCommand(message) {
  const grid = new P2P(null)
  grid.join('localhost', 3000, conn => {
    conn.on('data', data => {
      const answer = JSON.parse(data.toString('utf8'))
      console.log(answer)
      // TODO grid.closeConnection()
    })
    grid.write(conn, message)
  })
}

if (command === 'join') {
  BlockChain.mine('bar')
  const grid = new P2P(BlockChain)
  grid.join('localhost', 3000)
} else if (command === 'blockchain') {
  sendCommand(Messages.getBlockchain())
} else if (command === 'peers') {
  sendCommand(Messages.getPeers())
} else {
  const grid = new P2P(BlockChain)
  grid.serve()
}
