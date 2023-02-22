-- 000026_redelegation.up.sql
CREATE TABLE IF NOT EXISTS spacebox.redelegation_topic
(
    `delegator_address`     String,
    `src_validator_address` String,
    `dst_validator_address` String,
    `coin`                  String, -- TODO: @malekviktor this is json
    `height`                Int64,
    `completion_time`       TIMESTAMP
) ENGINE = Kafka('kafka:9093', 'redelegation', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.redelegation
(
    `delegator_address`     String,
    `src_validator_address` String,
    `dst_validator_address` String,
    `coin`                  String, -- TODO: @malekviktor this is json
    `height`                Int64,
    `completion_time`       TIMESTAMP
) ENGINE = ReplacingMergeTree()
      ORDER BY (`delegator_address`, `src_validator_address`, `dst_validator_address`, `height`, `completion_time`);

CREATE MATERIALIZED VIEW IF NOT EXISTS redelegation_consumer TO spacebox.redelegation
AS
SELECT delegator_address, src_validator_address, dst_validator_address, coin, height, completion_time
FROM spacebox.redelegation_topic
GROUP BY delegator_address, src_validator_address, dst_validator_address, coin, height, completion_time;
