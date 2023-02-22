-- 000015_supply.up.sql
CREATE TABLE IF NOT EXISTS spacebox.supply_topic
(
    `coins`  String, -- TODO: @malekviktor this is json
    `height` Int64
) ENGINE = Kafka('kafka:9093', 'supply', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.supply
(
    `coins`  String, -- TODO: @malekviktor this is json
    `height` Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`height`);

CREATE MATERIALIZED VIEW IF NOT EXISTS supply_consumer TO spacebox.supply
AS
SELECT coins, height
FROM spacebox.supply_topic
GROUP BY coins, height;
