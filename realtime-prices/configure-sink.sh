#!/usr/bin/env bash

curl \
    -X POST \
    -H "Accept:application/json" \
    -H "Content-Type:application/json" \
    --data @postgres-sink-btc.json http://localhost:8083/connectors