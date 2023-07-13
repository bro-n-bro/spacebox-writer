-- 000039_distribution_reward.up.sql
CREATE TABLE IF NOT EXISTS spacebox.distribution_reward_topic
(
    message String
) ENGINE = Kafka('kafka:9093', 'distribution_reward', 'spacebox', 'JSONAsString');


CREATE TABLE IF NOT EXISTS spacebox.distribution_reward
(
    `validator` String,
    `amount`    String,
    `height`    Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`height`, `validator`);

CREATE MATERIALIZED VIEW IF NOT EXISTS distribution_reward_consumer TO spacebox.distribution_reward
AS
SELECT JSONExtractString(message, 'validator') as validator,
       JSONExtractString(message, 'amount')    as amount,
       JSONExtractInt(message, 'height')       as height
FROM spacebox.distribution_reward_topic
GROUP BY validator, amount, height;

