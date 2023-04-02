-- 000030_delegation.up.sql
CREATE TABLE IF NOT EXISTS spacebox.delegation
(
    `operator_address`  String,
    `delegator_address` String,
    `coin`              json,
    `height`            Int64
) ENGINE = ReplacingMergeTree(`height`)
      ORDER BY (`operator_address`, `delegator_address`);
