const web3 = require('./connect')
const method = require('./method')

console.log(`fetching logs of contract ${method.contract}::${method.signature}`)
// retrieve logs, filtered by a specific method/event of a specific contract.
web3.eth.getPastLogs({
    // pulling the logs of the contract deployed at this address
    address: method.contract,
    // in the query as in the response, topics[0] is the method signature hash
    topics: [method.hash]
}).then(logs => {
    const transfers = logs.map(eventLog => {
        console.log(`[block ${eventLog.blockNumber}] decoding log #${eventLog.logIndex}: ${eventLog.id}`)
        decoded = web3.eth.abi.decodeLog(method.inputs, eventLog.data, eventLog.topics)
        console.log(`event: transfer ${decoded.value} from ${decoded.from} to ${decoded.to}`)

        return {
            txn: eventLog.transactionHash,
            contract: eventLog.address,
            event: {
                method: method.signature,
                signature: eventLog.topics[0],
            },
            params: {
                from: decoded.from,
                to: decoded.to,
                value: decoded.value,
            }
        }
    })

    console.log(transfers)
}).catch(console.error)