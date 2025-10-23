# [Analyze cryptocurrency market data](https://docs.timescale.com/timescaledb/latest/tutorials/analyze-cryptocurrency-data/#analyze-cryptocurrency-market-data)

## Installation

Assuming Docker and `docker-compose` installed, spin up TimescaleDB and Grafana: `docker compise up -d`.

Then setup environment variables and Python virtualenv:

```sh
cp .env.sample .env.dev
# edit .env.dev if necessary

source etc/up.sh`.
```