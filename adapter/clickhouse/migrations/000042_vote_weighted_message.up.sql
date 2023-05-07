-- 000042_vote_weighted_message.up.sql

CREATE TABLE IF NOT EXISTS spacebox.vote_weighted_message
(
    `height`               Int64,
    `msg_index`            Int64,
    `proposal_id`          Int64,
    `tx_hash`              String,
    `voter`                String,
    `weighted_vote_option` json
) ENGINE = ReplacingMergeTree()
      ORDER BY (`tx_hash`, `msg_index`);
