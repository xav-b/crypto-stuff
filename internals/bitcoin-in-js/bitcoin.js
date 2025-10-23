/**
 * Simple implementation of bitcoin network on top of the blockchain.
 *
 * https://medium.freecodecamp.org/how-does-bitcoin-work-i-built-an-app-to-show-you-f9fcd50bdd0d
 *
 */

const crypto = require('crypto')

const blockchain = require('./blockchain')

// FIXME Is transaction a block ??
// Transactions are a record of payment between two parties. When there is an
// exchange of value, a transaction is created to record it
class Transaction {
  constructor(type, inputs, outputs) {
    // can be:
    // - Reward: Satoshi rewarded with 100 coins for mining new block
    // - Regular: Satoshi paid Dean 5 coins with change of 94 coins
    // - Fee: Mining fee of 1 for whoever mines the transaction (Satoshi in example above)
    this.type = type
    // Where value is coming from
    this.inputs = inputs
    // Where value is going to
    this.outputs = outputs
  }

  // should it be part of blockchain ?
  hash() {
    // Uniquely identifies the transaction (using inputs & outputs)
    // TODO to be computed: f(this.inputs + this.outputs) = 000abcdefg…
    const algo = crypto.createHash('sha256')
    return algo
      .update(JSON.stringify(this.inputs))
      .update(JSON.stringify(this.outputs))
      .digest('hex')
  }
}

// FIXME I think it should extend `Block` (given `Input`)
class Output {
  constructor(address, amount) {
    // What is the public wallet address to send the coins to?
    this.address = address
    // How many coins
    this.amount = amount
  }
}

// FIXME is Input a block ? And is outputBlock actually outputTransaction ?
// (which itself would be a block ?)
class Input {
  constructor(outputBlock, signature) {
    // Transaction hash of the (unspent) output
    this.transactionHash = outputBlock.hash
    // The index of the (unspent) output in the transaction
    this.outputIndex = outputBlock.index
    // Amount of the (unspent) output
    this.amount = outputBlock.amount
    // Address of the (unspent) output
    this.address = outputBlock.address
    // Signed by the Address’s private key
    // TODO where does it come from
    this.signature = signature
  }
}

function rewardTransaction(address, amount) {
  // outpurs created from mining a block
  const outputs = [new Output(address, amount)]
  // Reward transactions are created as a result of finding a valid block on
  // the blockchain. As a result, reward transactions do not have any inputs
  // because it creates new coins
  return new Transaction('reward', [], outputs)
}

function getBalance(address) {
  const inputs = getUnspentInputs()
  const inputsForAddress = inputs.filter(input => input.address === address)
  return inputsForAddress.reduce((total, input) => total + input.amount, 0)
}

function createTransaction(myAddress) {
  const paidTotal = payments.reduce((total, paid) => total + paid.amount, 0)
  const fee = payments.reduce((total, payment) => total + payment.fee, 0)
  const unspentInputs = getUnspentInputs()

  let inputTotal = 0
  const inputs = unspentInputs.takeUntil(output => {
    inputTotal = inputTotal + output.amount
    return inputTotal >= paidTotal + fee
  })

  let outputs = payments.map(payment => ({
    amount: payment.amount,
    address: payment.address,
  }))
  let change = inputTotal - paidTotal - fee
  if (change > 0) {
    const changeOutput = { amount: change, address: myAddress }
    outputs = outputs.push(changeOutput)
  }

  const signedInputs = signInputs(inputs, password)
  return new Transaction('regular', inputs, outputs)
}

function isTransactionDoubleSpent(allTransactions, transaction) {
  const allTransactions = getTransactions()

  return allTransactions.some(tx => tx.inputs.some(input => transaction.hasSameInput(input)))
}
