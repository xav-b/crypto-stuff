# Streaming relatime crypto prices

**Inspired from**:
- [Using Kafka Streams to Analyze Live Trading Activity for Crypto Exchanges](https://www.confluent.io/kafka-summit-lon19/using-kafka-streams-analyze-trading-crypto-exchanges/)
- [Realtime Crypto Tracker with Kafka and QuestDB](https://medium.com/swlh/realtime-crypto-tracker-with-kafka-and-questdb-b33b19048fc2)
- [Yitaek/kafka-crypto-questdb](https://github.com/Yitaek/kafka-crypto-questdb)

## Usage

```console
alias dco='docker compose'

dco up -d

# ... wait a bit ... verify everything started
dco ps

./configure-sink.sh
```