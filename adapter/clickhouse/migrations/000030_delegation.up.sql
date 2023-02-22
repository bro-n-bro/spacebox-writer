-- 000030_delegation.up.sql
CREATE TABLE IF NOT EXISTS spacebox.delegation_topic
(
    `operator_address`  String,
    `delegator_address` String,
    `coin`              String, -- TODO: @malekviktor this is json
    `height`            Int64
) ENGINE = Kafka('kafka:9093', 'delegation', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.delegation
(
    `operator_address`  String,
    `delegator_address` String,
    `coin`              String, -- TODO: @malekviktor this is json
    `height`            Int64
) ENGINE = ReplacingMergeTree(`height`)
      ORDER BY (`operator_address`, `delegator_address`);

CREATE MATERIALIZED VIEW IF NOT EXISTS delegation_consumer TO spacebox.delegation
AS
SELECT operator_address, delegator_address, coin, height
FROM spacebox.delegation_topic
GROUP BY operator_address, delegator_address, coin, height;