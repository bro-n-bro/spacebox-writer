-- 000039_distribution_reward.up.sql
CREATE TABLE IF NOT EXISTS spacebox.distribution_reward
(
    `height`       Int64,
    `validator`    String,
    `amount`       json
) ENGINE = ReplacingMergeTree()
      ORDER BY (`tx_hash`, `msg_index`);