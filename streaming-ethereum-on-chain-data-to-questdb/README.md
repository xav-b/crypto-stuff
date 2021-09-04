# [Streaming Ethereum On-Chain Data to QuestDB](https://medium.com/geekculture/streaming-ethereum-on-chain-data-to-questdb-ea6b51d990ab)

Credits to Yitaek Hwang, awesome content. Also available on [QuestDB blog](https://questdb.io/tutorial/2021/04/12/ethereum/).

Other inspiration resources:
- [Go Ethereum book](https://goethereumbook.org)
- [Eth JSON-RPC API](https://eth.wiki/json-rpc/API)

**On-chain data streamed**:

- [ ] blocks
- [ ] tokens
- [ ] token transfers
- [ ] transactions
- [ ] logs
- [ ] traces
- [ ] contracts


## Schemas

TODO: make them sql files and automate the process.

Stolen from [etl blockchain](https://github.com/blockchain-etl/ethereum-etl-postgres/tree/master/schema), with types converted to java for QuestDB.

```sql
create table block
(
    created_at timestamp,
    timestamp timestamp,

    number long,
    hash string,
    parent_hash string,
    nonce long256,
    sha3_uncles string,
    logs_bloom string,
    transactions_root string,
    state_root string,
    receipts_root string,
    miner string,
    difficulty long,
    total_difficulty long,
    size float,
    extra_data string,
    gas_limit int,
    gas_used int,
    transaction_count int,
    base_fee_per_gas long
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