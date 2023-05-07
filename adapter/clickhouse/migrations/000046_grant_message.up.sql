-- 000046_grant_message.up.sql
CREATE TABLE IF NOT EXISTS spacebox.grant_message
(
    `height`     Int64,
    `msg_index`  Int64,
    `tx_hash`    String,
    `granter`    String,
    `grantee`    String,
    `msg_type`   String,
    `expiration` TIMESTAMP
) ENGINE = ReplacingMergeTree()
      ORDER BY (`tx_hash`, `msg_index`);