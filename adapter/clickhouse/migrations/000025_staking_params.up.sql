-- 000025_staking_params.up.sql
CREATE TABLE IF NOT EXISTS spacebox.staking_params_topic
(
    `params` String, -- TODO: @malekviktor this is json
    `height` Int64
) ENGINE = Kafka('kafka:9093', 'staking_params', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.staking_params
(
    `params` String, -- TODO: @malekviktor this is json
    `height` Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`height`);

CREATE MATERIALIZED VIEW IF NOT EXISTS staking_params_consumer TO spacebox.staking_params
AS
SELECT params, height
FROM spacebox.staking_params_topic
GROUP BY params, height;
