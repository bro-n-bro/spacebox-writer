-- 000007_distribution_params.up.sql
CREATE TABLE IF NOT EXISTS spacebox.distribution_params
(
    `height` Int64,
    `params` json
) ENGINE = ReplacingMergeTree()
      ORDER BY (`height`);
