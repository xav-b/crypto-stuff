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
    number BIGINT NOT NULL,
    hash TEXT NOT NULL,
    parent_hash TEXT,
    -- nonce is too big for bigint...
    nonce TEXT,
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