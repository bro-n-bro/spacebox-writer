-- 000004_send_message.up.sql
CREATE TABLE IF NOT EXISTS spacebox.send_message
(
    `height`       Int64,
    `address_from` String,
    `address_to`   String,
    `tx_hash`      String,
    `coins`        json,
    `msg_index`    Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`tx_hash`, `msg_index`);