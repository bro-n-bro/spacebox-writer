CREATE TABLE my_test.account
(
    `address` String,
    `height`  UInt256
) ENGINE = MergeTree()
      PRIMARY KEY (address);

CREATE TABLE my_test.block
(
    `height`           UInt256,
    `hash`             String,
    `num_txs`          Int256,
    `total_gas`        UInt256,
    `proposer_address` String,
    `timestamp`        TIMESTAMP
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE my_test.message
(
    `transaction_hash` String,
    `index`            UInt256,
    `type`             String,
    `value`            json
) ENGINE = MergeTree()
      PRIMARY KEY (transaction_hash);

CREATE TABLE my_test.supply
(
    `height` UInt256,
    `coins`  json
) ENGINE = MergeTree()
      PRIMARY KEY (height);