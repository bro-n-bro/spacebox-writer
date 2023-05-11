-- 000050_handle_validator_signature.up.sql
CREATE TABLE IF NOT EXISTS spacebox.handle_validator_signature
(
    `height`  Int64,
    `address` String,
    `power`   String,
    `reason`  String,
    `jailed`  String,
    `burned`    json
) ENGINE = ReplacingMergeTree()
      ORDER BY (`height`, `address`);