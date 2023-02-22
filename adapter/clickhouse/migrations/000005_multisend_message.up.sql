-- 000005_multisend_message.up.sql
CREATE TABLE IF NOT EXISTS spacebox.multisend_message_topic
(
    `height`       Int64,
    `address_from` String,
    `addresses_to` Array(String),
    `tx_hash`      String,
    `coins`        String, -- TODO: @malekviktor this is JSON
    `msg_index`    Int64
) ENGINE = Kafka('kafka:9093', 'multisend_message', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.multisend_message
(
    `height`       Int64,
    `address_from` String,
    `addresses_to` Array(String),
    `tx_hash`      String,
    `coins`        String, -- TODO: @malekviktor this is JSON
    `msg_index`    Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`tx_hash`, `msg_index`);

CREATE MATERIALIZED VIEW IF NOT EXISTS multisend_message_consumer TO spacebox.multisend_message
AS
SELECT height, address_from, addresses_to, tx_hash, coins, msg_index
FROM spacebox.multisend_message_topic
GROUP BY height, address_from, addresses_to, tx_hash, coins, msg_index;
