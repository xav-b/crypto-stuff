# Streaming on-chain Ethereum data

## Getting started

**Pre-requesites**:

- Go installed (tested under version `go1.16.6 darwin/amd64`)
- Docker installed (tested under version `20.10.8`)
- You have an [Infura](https://infura.io/) account. Create a project and note the project ID.

Start a Postgres database and a Grafana dashboard (port 3000): `docker compose up -d`

```console
# Infura project ID for auth
export INFURA_PROJECT_ID="..."

# build the go binary `./stream`
make

# start listening for new blocks
./stream \
    -network mainnet \
    -db postgresql://postgres:RDLPWbx5hM3ra@localhost:5432/crypto

./stream \
    -network mainnet \
    -db postgresql://postgres:RDLPWbx5hM3ra@localhost:5432/crypto \
    -block 0xa5a821871e51e45437dc192321b080a9c0ced86aadefb9e3e0e071398fcc87a3
```

**Resources and credits**:
- [Go Ethereum book](https://goethereumbook.org)
- [Eth JSON-RPC API](https://eth.wiki/json-rpc/API)
- [Go-ethereum documentation](https://pkg.go.dev/github.com/ethereum/go-ethereum@v1.10.8)

**On-chain data streamed**:

- [x] blocks
- [ ] tokens
- [ ] token transfers
- [x] transactions
- [ ] logs
- [ ] traces
- [ ] contracts

**Fixmes**:

- [ ] Ctrl-c is no longer caught since I added transactions storage
- [ ] Some values get bigger than Postgres `BIGINT`
- [ ] Data bytes fail to be encoded as UTF-8 for Postgres storage

---

Originally inspired from [Streaming Ethereum On-Chain Data to QuestDB](https://medium.com/geekculture/streaming-ethereum-on-chain-data-to-questdb-ea6b51d990ab).

Credits to Yitaek Hwang, awesome content. Also available on [QuestDB blog](https://questdb.io/tutorial/2021/04/12/ethereum/).

## QuestDB Schemas

In theory QuestDB can be used as a Postgres-compatible replacement so passing
`-db postgresql://admin:quest@localhost:8812/qdb` should work to initialise the
connection. Then use the schemas below, although I'm no longer sure they map to
the golang types (feel free to open an issue).

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