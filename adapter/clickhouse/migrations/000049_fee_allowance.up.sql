-- 000049_fee_allowance.up.sql
CREATE TABLE IF NOT EXISTS spacebox.fee_allowance
(
    `height`     Int64,
    `granter`    String,
    `grantee`    String,
    `allowance`  String,
    `expiration` TIMESTAMP
) ENGINE = ReplacingMergeTree()
      ORDER BY (`granter`, `grantee`);