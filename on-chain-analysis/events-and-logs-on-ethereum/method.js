const web3 = require('./connect')

const inputs = [
    {
        type: 'string',
        name: 'from',
        // when a parameter is indexed, it lands in the `topics` array of the log
        //
        // I'm not sure what purpose it is supposed to serve but web3 allows us
        // to filter by topics, hence we can filter onn those addresses.
        // Effectively we could fetch events for transactions sent to a given
        // address
        indexed: true,
    }, {
        type: 'string',
        name: 'to',
        indexed: true,
    }, {
        type: 'uint256',
        name: 'value',
        indexed: false,
    }
]
const signature = 'Transfer(address,address,uint256)'
// <=> '0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef'


module.exports = {
    // https://etherscan.io/tx/0xdadd97094b8a789d387a53845449d52213e9ffd37a3530d8d99d2234dea820fa#eventlog
    contract: '0xdac17f958d2ee523a2206206994597c13d831ec7',
    signature,
    hash: web3.utils.keccak256(signature),
    inputs,
}