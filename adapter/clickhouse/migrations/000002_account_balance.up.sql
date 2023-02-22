-- 000002_account_balance.up.sql
CREATE TABLE IF NOT EXISTS spacebox.account_balance_topic
(
    `address` String,
    `height`  Int64,
    `coins`   String -- TODO: @malekviktor this is JSON
) ENGINE = Kafka('kafka:9093', 'account_balance', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.account_balance
(
    `address` String,
    `height`  Int64,
    `coins`   String -- TODO: @malekviktor this is JSON

) ENGINE = ReplacingMergeTree(`height`)
      ORDER BY (`address`);

CREATE MATERIALIZED VIEW IF NOT EXISTS account_balance_consumer TO spacebox.account_balance
AS
SELECT address, height, coins
FROM spacebox.account_balance_topic
GROUP BY address, height, coins;
