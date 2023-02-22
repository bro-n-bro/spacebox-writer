-- 000024_mint_params.up.sql
CREATE TABLE IF NOT EXISTS spacebox.mint_params_topic
(
    `params` String, -- TODO: @malekviktor this is json
    `height` Int64
) ENGINE = Kafka('kafka:9093', 'mint_params', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.mint_params
(
    `params` String, -- TODO: @malekviktor this is json
    `height` Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`height`);

CREATE MATERIALIZED VIEW IF NOT EXISTS mint_params_consumer TO spacebox.mint_params
AS
SELECT params, height
FROM spacebox.mint_params_topic
GROUP BY params, height;
