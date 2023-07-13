-- 000028_unbonding_delegation.up.sql
CREATE TABLE IF NOT EXISTS spacebox.unbonding_delegation_topic
(
    message String
) ENGINE = Kafka('kafka:9093', 'unbonding_delegation', 'spacebox', 'JSONAsString');

CREATE TABLE IF NOT EXISTS spacebox.unbonding_delegation
(
    `validator_address`    String,
    `delegator_address`    String,
    `coin`                 String,
    `completion_timestamp` TIMESTAMP,
    `height`               Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`validator_address`, `delegator_address`, `completion_timestamp`, `height`);

CREATE MATERIALIZED VIEW IF NOT EXISTS unbonding_delegation_consumer TO spacebox.unbonding_delegation
AS
SELECT JSONExtractString(message, 'validator_address')                      as validator_address,
       JSONExtractString(message, 'delegator_address')                      as delegator_address,
       JSONExtractString(message, 'coin')                                   as coin,
       toDateTimeOrZero(JSONExtractString(message, 'completion_timestamp')) as completion_timestamp,
       JSONExtractInt(message, 'height')                                    as height
FROM spacebox.unbonding_delegation_topic
GROUP BY validator_address, delegator_address, coin, completion_timestamp, height;