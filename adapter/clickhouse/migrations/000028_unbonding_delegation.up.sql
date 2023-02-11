-- 000028_unbonding_delegation.up.sql
CREATE TABLE IF NOT EXISTS spacebox.unbonding_delegation
(
    `validator_address`    String,
    `delegator_address`    String,
    `coin`                 json,
    `completion_timestamp` TIMESTAMP,
    `height`               Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`validator_address`, `delegator_address`, `completion_timestamp`, `height`);
