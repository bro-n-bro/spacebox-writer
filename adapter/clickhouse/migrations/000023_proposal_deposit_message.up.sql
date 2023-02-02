-- 000023_proposal_deposit_message.up.sql
CREATE TABLE IF NOT EXISTS spacebox.proposal_deposit_message
(
    `proposal_id`       Int64,
    `depositor_address` String,
    `coins`             json,
    `height`            Int64,
    `tx_hash`           String,
    `msg_index`         Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`tx_hash`, `msg_index`);