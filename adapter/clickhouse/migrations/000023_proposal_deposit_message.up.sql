-- 000023_proposal_deposit_message.up.sql
CREATE TABLE IF NOT EXISTS spacebox.proposal_deposit_message_topic
(
    `proposal_id`       Int64,
    `depositor_address` String,
    `coins`             String, -- TODO: @malekviktor this is json
    `height`            Int64,
    `tx_hash`           String,
    `msg_index`         Int64
) ENGINE = Kafka('kafka:9093', 'proposal_deposit_message', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.proposal_deposit_message
(
    `proposal_id`       Int64,
    `depositor_address` String,
    `coins`             String, -- TODO: @malekviktor this is json
    `height`            Int64,
    `tx_hash`           String,
    `msg_index`         Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`tx_hash`, `msg_index`);


CREATE MATERIALIZED VIEW IF NOT EXISTS proposal_deposit_message_consumer TO spacebox.proposal_deposit_message
AS
SELECT proposal_id, depositor_address, coins, height, tx_hash, msg_index
FROM spacebox.proposal_deposit_message_topic
GROUP BY proposal_id, depositor_address, coins, height, tx_hash, msg_index;
