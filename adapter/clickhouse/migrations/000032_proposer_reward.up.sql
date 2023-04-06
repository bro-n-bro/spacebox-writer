-- 000032_proposer_reward.up.sql
CREATE TABLE IF NOT EXISTS spacebox.proposer_reward
(
    `height`    Int64,
    `validator` String,
    `reward`    json
) ENGINE = ReplacingMergeTree()
      ORDER BY (`height`, `validator`);
