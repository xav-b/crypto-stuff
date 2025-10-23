/**
 * Following this tutorial: https://www.trustology.io/insights-events/decoding-an-ethereum-transaction-its-no-secret-its-just-smart
 * 
 * 
 */

const web3 = require('./connect')

// Hash example: https://etherscan.io/tx/0xdadd97094b8a789d387a53845449d52213e9ffd37a3530d8d99d2234dea820fa
const txnHash = '0xdadd97094b8a789d387a53845449d52213e9ffd37a3530d8d99d2234dea820fa'
const method = 'execTransaction(address,uint256,bytes,uint8,uint256,uint256,uint256,address,address,bytes)'

const methodHash = web3.utils.keccak256(method)
const methodId = methodHash.slice(0, 2 + 8)  // '0x' + 8 firt hex characters (4 bytes)
// should match 6a761202
// alternatively: web3.eth.abi.encodeFunctionSignature()
// interesting note: `a9059cbb` is the famous method ID of an ERC20 transfer: transfer(address,uint256)

console.log(`method hash: ${methodHash}\nmethod ID: ${methodId}`)

web3.eth.getTransaction(txnHash)
    .then(onTxn)
    .catch(console.error)

const removePadding = (str) => str.replace(/^0+/gm,'')

function onTxn(txn) {
    // we want to decode the transaction payload - or in other words the smart
    // contract signature and passed parameters

    // first extract the 0x and method id
    const txnMethodId = txn.input.slice(0, 10)
    console.log(`verifying method ID: ${txnMethodId === methodId}`)

    const payload = txn.input.slice(10)
    // break down the remaining data in 32 bytes (64 hex chars) blocks
    const fields = []
    // start our block pointer after the method id
    let pointer = 10
    while (pointer + 64 <= txn.input.length) {
       block = txn.input.slice(pointer, pointer + 64) 
       fields.push(block)
       console.log(pointer, block)
       pointer = pointer + 64
    }

    // now let's decode the blocks, following the method signature
    // after the method id, each block is a parameter, encoded but of the same
    // type as defined in the signature aboive

    // the first parameter is the address we are sending to
    // this happens here to be the Tether contract
    const toAddress = removePadding(fields[0])
    console.log(`param to: ${toAddress}`)

    // value being transacted in Ether
    // TODO: convert to uint256
    const value = fields[1]
    console.log(`param value: ${value}`)

    // next is the offset in hex of bytes from method id to the remaining data
    let offset = parseInt(removePadding(fields[2]), 16)  // in bytes
    offset = offset * 2  // in hex chars
    offset = offset / 64  // fields index
    console.log(`jumping at offset row #${offset}`)

    // ... next fields hold some information on gas, token ...

    const dataLength = parseInt(fields[10], 16) * 2
    console.log(`data is ${dataLength} hex chars long`)
    // note that dataLength is padded with 0s to reach a rounded block
    const dataPayload = txn.input.slice(10 + 11 * 64, 10+ 11 * 64 + dataLength)
    console.log('data payload:', dataPayload)
    // TODO: parse dataPayload, which is a method call. First bytes are the
    // method id, followed by an address and a value
}