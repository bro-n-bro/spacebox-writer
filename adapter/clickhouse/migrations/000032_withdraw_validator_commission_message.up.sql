-- 000032_withdraw_validator_commission_message.up.sql
CREATE TABLE IF NOT EXISTS spacebox.withdraw_validator_commission_message_topic
(
    `height`              Int64,
    `tx_hash`             String,
    `msg_index`           Int64,
    `withdraw_commission` Int64,
    `sender_address`      String
) ENGINE = Kafka('kafka:9093', 'withdraw_validator_commission_message', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.withdraw_validator_commission_message
(
    `height`              Int64,
    `tx_hash`             String,
    `msg_index`           Int64,
    `withdraw_commission` Int64,
    `sender_address`      String
) ENGINE = ReplacingMergeTree()
      ORDER BY (`tx_hash`, `msg_index`);

CREATE MATERIALIZED VIEW IF NOT EXISTS withdraw_validator_commission_message_consumer TO spacebox.withdraw_validator_commission_message
AS
SELECT height, tx_hash, msg_index, withdraw_commission, sender_address
FROM spacebox.withdraw_validator_commission_message_topic
GROUP BY height, tx_hash, msg_index, withdraw_commission, sender_address;

