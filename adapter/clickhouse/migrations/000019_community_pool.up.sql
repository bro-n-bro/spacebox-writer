-- 000019_community_pool.up.sql
CREATE TABLE IF NOT EXISTS spacebox.community_pool
(
    `coins`  json,
    `height` Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`height`);
