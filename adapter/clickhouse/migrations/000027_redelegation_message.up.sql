-- 000027_redelegation_message.up.sql
CREATE TABLE IF NOT EXISTS spacebox.redelegation_message_topic
(
    `delegator_address`     String,
    `src_validator_address` String,
    `dst_validator_address` String,
    `coin`                  String, -- TODO: @malekviktor this is json
    `height`                Int64,
    `completion_time`       TIMESTAMP,
    `tx_hash`               String,
    `msg_index`             Int64
) ENGINE = Kafka('kafka:9093', 'redelegation_message', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.redelegation_message
(
    `delegator_address`     String,
    `src_validator_address` String,
    `dst_validator_address` String,
    `coin`                  String, -- TODO: @malekviktor this is json
    `height`                Int64,
    `completion_time`       TIMESTAMP,
    `tx_hash`               String,
    `msg_index`             Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`tx_hash`, `msg_index`);

CREATE MATERIALIZED VIEW IF NOT EXISTS redelegation_message_consumer TO spacebox.redelegation_message
AS
SELECT delegator_address,
       src_validator_address,
       dst_validator_address,
       coin,
       height,
       completion_time,
       tx_hash,
       msg_index
FROM spacebox.redelegation_message_topic
GROUP BY delegator_address, src_validator_address, dst_validator_address, coin, height, completion_time, tx_hash,
         msg_index;