# [Streaming Ethereum On-Chain Data to QuestDB](https://medium.com/geekculture/streaming-ethereum-on-chain-data-to-questdb-ea6b51d990ab)

Credits to Yitaek Hwang, awesome content.


## Schemas

Stolen from [etl blockchain](https://github.com/blockchain-etl/ethereum-etl-postgres/tree/master/schema), with types converted to java for QuestDB.

```sql
create table blocks
(
    timestamp timestamp,

    number bigint,
    hash string,
    parent_hash string,
    nonce string,
    sha3_uncles string,
    logs_bloom string,
    transactions_root string,
    state_root string,
    receipts_root string,
    miner string,
    difficulty string,
    total_difficulty string,
    size bigint,
    extra_data string,
    gas_limit bigint,
    gas_used bigint,
    transaction_count bigint,
    base_fee_per_gas bigint
);
```

```sql
create table token_transfers
(
    token_address string,
    from_address string,
    to_address string,
    value long,
    transaction_hash string,
    log_index bigint,
    block_timestamp timestamp,
    block_number bigint,
    block_hash string
);
```

```sql
create table transactions
(
    hash string,
    nonce bigint,
    transaction_index bigint,
    from_address string,
    to_address string,
    value long,
    gas bigint,
    gas_price bigint,
    input string,
    receipt_cumulative_gas_used bigint,
    receipt_gas_used bigint,
    receipt_contract_address string,
    receipt_root string,
    receipt_status bigint,
    block_timestamp timestamp,
    block_number bigint,
    block_hash string,
    max_fee_per_gas bigint,
    max_priority_fee_per_gas bigint,
    transaction_type bigint,
    receipt_effective_gas_price bigint
);
```