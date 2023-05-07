-- 000041_edit_validator_message.up.sql

CREATE TABLE IF NOT EXISTS spacebox.edit_validator_message
(
    `height`       Int64,
    `msg_index`    Int64,
    `tx_hash`      String,
    `description`  String
) ENGINE = ReplacingMergeTree()
    ORDER BY (`tx_hash`, `msg_index`);