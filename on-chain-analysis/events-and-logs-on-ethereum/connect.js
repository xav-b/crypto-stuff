/**
 * Singleton to connect to Ethereum.
 */

const Web3 = require('web3')

// we use the very populat Infura service to avoid having to run a local node.
// Ethereum being a distributed blockchain, there's no central API to reach out.
const network = 'mainnet'
// const providerUri = `wss://${network}.infura.io/ws/v3/${process.env.INFURA_PROJECT_ID}`
const providerUri = `https://${network}.infura.io/v3/${process.env.INFURA_PROJECT_ID}`

console.log(`connecting to Ethereum using provider: ${providerUri}`)
// module.exports = new Web3(new Web3.providers.WebsocketProvider(providerUri))
module.exports = new Web3(new Web3.providers.HttpProvider(providerUri))