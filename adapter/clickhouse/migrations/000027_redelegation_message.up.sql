-- 000027_redelegation_message.up.sql
CREATE TABLE IF NOT EXISTS spacebox.redelegation_message
(
    `delegator_address`     String,
    `src_validator_address` String,
    `dst_validator_address` String,
    `coin`                  json,
    `height`                Int64,
    `completion_time`       TIMESTAMP,
    `tx_hash`               String,
    `msg_index`             Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`tx_hash`, `msg_index`);
