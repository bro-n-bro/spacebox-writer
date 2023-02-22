-- 000028_unbonding_delegation.up.sql
CREATE TABLE IF NOT EXISTS spacebox.unbonding_delegation_topic
(
    `validator_address`    String,
    `delegator_address`    String,
    `coin`                 String, -- TODO: @malekviktor this is json
    `completion_timestamp` TIMESTAMP,
    `height`               Int64
) ENGINE = Kafka('kafka:9093', 'unbonding_delegation', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.unbonding_delegation
(
    `validator_address`    String,
    `delegator_address`    String,
    `coin`                 String, -- TODO: @malekviktor this is json
    `completion_timestamp` TIMESTAMP,
    `height`               Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`validator_address`, `delegator_address`, `completion_timestamp`, `height`);


CREATE MATERIALIZED VIEW IF NOT EXISTS unbonding_delegation_consumer TO spacebox.unbonding_delegation
AS
SELECT validator_address, delegator_address, coin, completion_timestamp, height
FROM spacebox.unbonding_delegation_topic
GROUP BY validator_address, delegator_address, coin, completion_timestamp, height;