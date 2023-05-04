-- 000033_submit_proposal_message.up.sql
CREATE TABLE IF NOT EXISTS spacebox.submit_proposal_message
(
    `tx_hash`         String,
    `proposer`        String,
    `initial_deposit` json,
    `title`           String,
    `description`     String,
    `type`            String,
    `proposal_id`     Int64,
    `height`          Int64,
    `msg_index`       Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`tx_hash`, `msg_index`);
