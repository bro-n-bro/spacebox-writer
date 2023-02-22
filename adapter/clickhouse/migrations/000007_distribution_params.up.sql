-- 000007_distribution_params.up.sql
CREATE TABLE IF NOT EXISTS spacebox.distribution_params_topic
(
    `height` Int64,
    `params` String -- TODO: @malekviktor this is JSON
) ENGINE = Kafka('kafka:9093', 'distribution_params', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.distribution_params
(
    `height` Int64,
    `params` String -- TODO: @malekviktor this is JSON
) ENGINE = ReplacingMergeTree()
      ORDER BY (`height`);

CREATE MATERIALIZED VIEW IF NOT EXISTS distribution_params_consumer TO spacebox.distribution_params
AS
SELECT height, params
FROM spacebox.distribution_params_topic
GROUP BY height, params;
