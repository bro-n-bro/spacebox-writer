-- 000029_unbonding_delegation_message.up.sql
CREATE TABLE IF NOT EXISTS spacebox.unbonding_delegation_message
(
    `validator_address`    String,
    `delegator_address`    String,
    `coin`                 json,
    `completion_timestamp` TIMESTAMP,
    `height`               Int64,
    `tx_hash`              String,
    `msg_index`            Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`tx_hash`, `msg_index`)
