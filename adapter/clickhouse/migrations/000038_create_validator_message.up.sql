-- 000038_create_validator_message.up.sql

CREATE TABLE IF NOT EXISTS spacebox.create_validator_message_topic
(
    `height`              Int64,
    `msg_index`           Int64,
    `tx_hash`             String,
    `delegator_address`   String,
    `validator_address`   String,
    `description`         String,
    `commission_rates`    String,
    `min_self_delegation` Int64,
    `pubkey`              String
) ENGINE = Kafka('kafka:9093', 'create_validator_message', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.create_validator_message
(
    `height`              Int64,
    `msg_index`           Int64,
    `tx_hash`             String,
    `delegator_address`   String,
    `validator_address`   String,
    `description`         String,
    `commission_rates`    String,
    `min_self_delegation` Int64,
    `pubkey`              String
) ENGINE = ReplacingMergeTree()
      ORDER BY (`tx_hash`, `msg_index`);

CREATE MATERIALIZED VIEW IF NOT EXISTS create_validator_message_consumer TO spacebox.create_validator_message
AS
SELECT height,
       msg_index,
       tx_hash,
       delegator_address,
       validator_address,
       description,
       commission_rates,
       min_self_delegation,
       pubkey
FROM spacebox.create_validator_message_topic
GROUP BY height,
         msg_index,
         tx_hash,
         delegator_address,
         validator_address,
         description,
         commission_rates,
         min_self_delegation,
         pubkey;

