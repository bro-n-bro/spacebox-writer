-- 000020_delegation_reward_message.up.sql
CREATE TABLE IF NOT EXISTS spacebox.delegation_reward_message_topic
(
    `validator_address` String,
    `delegator_address` String,
    `coins`             String, -- TODO: @malekviktor this is json
    `height`            Int64,
    `tx_hash`           String,
    `msg_index`         Int64
) ENGINE = Kafka('kafka:9093', 'delegation_reward_message', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.delegation_reward_message
(
    `validator_address` String,
    `delegator_address` String,
    `coins`             String, -- TODO: @malekviktor this is json
    `height`            Int64,
    `tx_hash`           String,
    `msg_index`         Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`tx_hash`, `msg_index`);

CREATE MATERIALIZED VIEW IF NOT EXISTS delegation_reward_message_consumer TO spacebox.delegation_reward_message
AS
SELECT validator_address, delegator_address, coins, height, tx_hash, msg_index
FROM spacebox.delegation_reward_message_topic
GROUP BY validator_address, delegator_address, coins, height, tx_hash, msg_index;
