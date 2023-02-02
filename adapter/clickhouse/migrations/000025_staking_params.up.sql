-- 000025_staking_params.up.sql
CREATE TABLE IF NOT EXISTS spacebox.staking_params
(
    `params` json,
    `height` Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`height`);
