-- 000044_grant_allowance_message.up.sql
CREATE TABLE IF NOT EXISTS spacebox.grant_allowance_message
(
    `height`     Int64,
    `msg_index`  Int64,
    `tx_hash`    String,
    `granter`    String,
    `grantee`    String,
    `allowance`  String,
    `expiration` TIMESTAMP
) ENGINE = ReplacingMergeTree()
      ORDER BY (`tx_hash`, `msg_index`);