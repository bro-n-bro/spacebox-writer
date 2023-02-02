-- 000005_multisend_message.up.sql
CREATE TABLE IF NOT EXISTS spacebox.multisend_message
(
    `height`       Int64,
    `address_from` String,
    `addresses_to` Array(String),
    `tx_hash`      String,
    `coins`        json,
    `msg_index`    Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`tx_hash`, `msg_index`);