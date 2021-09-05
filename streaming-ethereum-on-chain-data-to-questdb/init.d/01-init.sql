/*
 * Initialise PostgreSQL
 */

 CREATE DATABASE crypto;

 \connect crypto;

 CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
 CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE;

-- TODO: use fixed-char for known string length like addresses and hashes
CREATE TABLE IF NOT EXISTS block (
    -- FIXME: I think timezone should be supported
    ts TIMESTAMP NOT NULL,
    -- TODO: rename block_bumber and block_hash
    number BIGINT NOT NULL,
    hash TEXT NOT NULL,
    parent_hash TEXT,
    -- nonce is too big for bigint...
    nonce TEXT NOT NULL,

    sha3_uncles TEXT,
    logs_bloom TEXT,
    transactions_root TEXT,
    state_root TEXT,
    receipts_root TEXT,
    miner TEXT,
    difficulty BIGINT,
    total_difficulty BIGINT,
    size REAL,
    extra_data TEXT,
    gas_limit INTEGER,
    gas_used INTEGER,
    transaction_count INTEGER,
    base_fee_per_gas BIGINT,

    created_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS tx (
    tx_hash TEXT NOT NULL,
    nonce TEXT NOT NULL,
    tx_index INTEGER NOT NULL,
    tx_type BIGINT,
    from_address TEXT NOT NULL,
    to_address TEXT NOT NULL,
    input TEXT,

    -- prices
    tx_value BIGINT,
    -- NOTE Deprecated tx_cost BIGINT, -- use gas*gasPrice + value
    gas BIGINT,
    gas_price BIGINT,
    max_fee_per_gas BIGINT,
    max_priority_fee_per_gas BIGINT,

    -- receipt
    receipt_cumulative_gas_used BIGINT,
    receipt_gas_used BIGINT,
    receipt_contract_address TEXT,
    receipt_root TEXT,
    receipt_status BIGINT,
    receipt_effective_gas_price BIGINT,

    -- block
    block_number BIGINT,
    block_hash TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT now()
);