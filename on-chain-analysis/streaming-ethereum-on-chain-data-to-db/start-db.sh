#! /usr/bin/env bash

# More info: https://questdb.io/docs/get-started/docker/

# unofficial strict mode
set -eo pipefail

# port 8812: Postgres wire protocol
# port 9000: REST API and Web Console
docker run \
  --detach --name cryptodb \
  -v "$(pwd)/_db:/root/.questdb/" \
  -p 9000:9000 \
  -p 8812:8812 \
  questdb/questdb
