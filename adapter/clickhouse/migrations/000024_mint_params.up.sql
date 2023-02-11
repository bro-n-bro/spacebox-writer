-- 000024_mint_params.up.sql
CREATE TABLE IF NOT EXISTS spacebox.mint_params
(
    `params` json,
    `height` Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`height`);
