# Spacebox-writer

![linter](https://github.com/bro-n-bro/spacebox-writer/actions/workflows/lint.yml/badge.svg)

Spacebox-writer is another one of the crucial parts of [Spacebox](https://github.com/bro-n-bro/spacebox) indexer. The writer pulls specific topics from Kafka broker and writes them to Clickhouse DB. Simultaneously it self-logs its progress to Mongo-db to ensure consistency.

## Installation

### From source

```bash
git clone
cd spacebox-writer
go build
```

### From binary

Download the latest release from the [releases page](...)

## Usage

```bash
./spacebox-writer
```

## Configuration

The configuration file is ENV based. You can set the following environment variables:

**// fixme: update this table for the new env vars**

| Variable                   | Description                                  | Default     |
|----------------------------|----------------------------------------------|-------------|
| `CLICKHOUSE_HOST`          | The host of the clickhouse database          | `localhost` |
| `CLICKHOUSE_PORT`          | The port of the clickhouse database          | `9000`      |
| `CLICKHOUSE_USER`          | The user of the clickhouse database          | `default`   |
| `CLICKHOUSE_PASSWORD`      | The password of the clickhouse database      | `""`        |
| `CLICKHOUSE_DATABASE`      | The database of the clickhouse database      | `default`   |
| `CLICKHOUSE_TABLE`         | The table of the clickhouse database         | `default`   |
| `CLICKHOUSE_TIMEOUT`       | The timeout of the clickhouse database       | `10`        |
| `CLICKHOUSE_MAX_RETRIES`   | The max retries of the clickhouse database   | `3`         |
| `CLICKHOUSE_RETRY_DELAY`   | The retry delay of the clickhouse database   | `3`         |
| `CLICKHOUSE_BATCH_SIZE`    | The batch size of the clickhouse database    | `1000`      |
| `CLICKHOUSE_BATCH_TIMEOUT` | The batch timeout of the clickhouse database | `1`         |

## Docker

### Build

```bash
docker build -t spacebox-writer .
```

### Run

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

### Before commiting

- run `make fmt` to format the code.
- run `make test` to check for test errors.
- run `make lint` to check for linting errors.
- run `make fix` to fix struct sizes.
