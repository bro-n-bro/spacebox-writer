-- 000019_community_pool.up.sql
CREATE TABLE IF NOT EXISTS spacebox.community_pool_topic
(
    `coins`  String, -- TODO: @malekviktor this is json
    `height` Int64
) ENGINE = Kafka('kafka:9093', 'community_pool', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.community_pool
(
    `coins`  String, -- TODO: @malekviktor this is json
    `height` Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`height`);

CREATE MATERIALIZED VIEW IF NOT EXISTS community_pool_consumer TO spacebox.community_pool
AS
SELECT coins, height
FROM spacebox.community_pool_topic
GROUP BY coins, height;
