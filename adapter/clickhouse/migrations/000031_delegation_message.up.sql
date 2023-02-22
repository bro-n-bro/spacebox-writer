-- 000031_delegation_message.up.sql
CREATE TABLE IF NOT EXISTS spacebox.delegation_message_topic
(
    `operator_address`  String,
    `delegator_address` String,
    `coin`              String, -- TODO: @malekviktor this is json
    `height`            Int64,
    `tx_hash`           String,
    `msg_index`         Int64
) ENGINE = Kafka('kafka:9093', 'delegation_message', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.delegation_message
(
    `operator_address`  String,
    `delegator_address` String,
    `coin`              String, -- TODO: @malekviktor this is json
    `height`            Int64,
    `tx_hash`           String,
    `msg_index`         Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`tx_hash`, `msg_index`);

CREATE MATERIALIZED VIEW IF NOT EXISTS delegation_message_consumer TO spacebox.delegation_message
AS
SELECT operator_address, delegator_address, coin, height, tx_hash, msg_index
FROM spacebox.delegation_message_topic
GROUP BY operator_address, delegator_address, coin, height, tx_hash, msg_index;