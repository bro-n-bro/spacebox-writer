-- 000059_validator_voting_power.up.sql
CREATE TABLE IF NOT EXISTS spacebox.validator_voting_power_topic
(
    `height`            Int64,
    `voting_power`      Int64,
    `validator_address` String,
    `timestamp`         String
) ENGINE = Kafka('kafka:9093', 'validator_voting_power', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.validator_voting_power
(
    `height`            Int64,
    `voting_power`      Int64,
    `validator_address` String,
    `timestamp`         TIMESTAMP
) ENGINE = ReplacingMergeTree()
      ORDER BY (`height`, `validator_address`);

CREATE MATERIALIZED VIEW IF NOT EXISTS validator_voting_power_consumer TO spacebox.validator_voting_power
AS
SELECT height,
       voting_power,
       validator_address,
       parseDateTimeBestEffortOrZero(timestamp) AS timestamp
FROM spacebox.validator_voting_power_topic
GROUP BY height, voting_power, validator_address, timestamp;
