# Spacebox-writer

Spacebox-writer is another one of the crucial parts of [Spacebox](https://github.com/bro-n-bro/spacebox) indexer. The writer pulls specific topics from Kafka broker and writes them to Clickhouse DB. Simultaneously it self-logs its progress to Mongo-db to ensure consistency.

## build

```bash
docker build -t spacebox-writer:latest .
```

## run

Running writer standalone is pretty much pointless, so please refer to the main [Sacebox repo](https://github.com/bro-n-bro/spacebox#readme) to find out how to start the whole setup.
