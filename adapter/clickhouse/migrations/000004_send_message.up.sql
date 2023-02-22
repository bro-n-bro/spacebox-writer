-- 000004_send_message.up.sql
CREATE TABLE IF NOT EXISTS spacebox.send_message_topic
(
    `height`       Int64,
    `address_from` String,
    `address_to`   String,
    `tx_hash`      String,
    `coins`        String, -- TODO: @malekviktor this is JSON
    `msg_index`    Int64
) ENGINE = Kafka('kafka:9093', 'send_message', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.send_message
(
    `height`       Int64,
    `address_from` String,
    `address_to`   String,
    `tx_hash`      String,
    `coins`        String, -- TODO: @malekviktor this is JSON
    `msg_index`    Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`tx_hash`, `msg_index`);

CREATE MATERIALIZED VIEW IF NOT EXISTS send_message_consumer TO spacebox.send_message
AS
SELECT height, address_from, address_to, tx_hash, coins, msg_index
FROM spacebox.send_message_topic
GROUP BY height, address_from, address_to, tx_hash, coins, msg_index;
