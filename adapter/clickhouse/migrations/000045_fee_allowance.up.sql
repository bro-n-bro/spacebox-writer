-- 000045_fee_allowance.up.sql
CREATE TABLE IF NOT EXISTS spacebox.fee_allowance
(
    `height`     Int64,
    `granter`    String,
    `grantee`    String,
    `allowance`  String,
    `expiration` TIMESTAMP,
    `is_active`  BOOLEAN
) ENGINE = ReplacingMergeTree()
      ORDER BY (`granter`, `grantee`);