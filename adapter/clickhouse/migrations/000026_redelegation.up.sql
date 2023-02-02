-- 000026_redelegation.up.sql
CREATE TABLE IF NOT EXISTS spacebox.redelegation
(
    `delegator_address`     String,
    `src_validator_address` String,
    `dst_validator_address` String,
    `coin`                  json,
    `height`                Int64,
    `completion_time`       TIMESTAMP
) ENGINE = ReplacingMergeTree()
      ORDER BY (`delegator_address`, `src_validator_address`, `dst_validator_address`, `height`, `completion_time`);
