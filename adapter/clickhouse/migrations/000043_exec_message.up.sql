-- 000043_exec_message.up.sql
CREATE TABLE IF NOT EXISTS spacebox.exec_message
(
    `height`    Int64,
    `msg_index` Int64,
    `tx_hash`   String,
    `grantee`   String,
    `msgs`      String
) ENGINE = ReplacingMergeTree()
      ORDER BY (`tx_hash`, `msg_index`);