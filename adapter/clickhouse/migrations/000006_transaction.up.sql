-- 000006_transaction.up.sql
CREATE TABLE IF NOT EXISTS spacebox.transaction_topic
(
    `signatures` Array(String),
    `hash`         String,
    `height`       Int64,
    `success`      BOOL,
    `messages`     String,
    `memo`         String,
    `signer_infos` String, -- TODO: @malekviktor this is JSON
    `fee`          String, -- TODO: @malekviktor this is JSON
    `signer`       String,
    `gas_wanted`   Int64,
    `gas_used`     Int64,
    `raw_log`      String,
    `logs`         String  -- TODO: @malekviktor this is JSON
) ENGINE = Kafka('kafka:9093', 'transaction', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.transaction
(
    `signatures` Array(String),
    `hash`         String,
    `height`       Int64,
    `success`      BOOL,
    `messages`     String,
    `memo`         String,
    `signer_infos` String, -- TODO: @malekviktor this is JSON
    `fee`          String, -- TODO: @malekviktor this is JSON
    `signer`       String,
    `gas_wanted`   Int64,
    `gas_used`     Int64,
    `raw_log`      String,
    `logs`         String  -- TODO: @malekviktor this is JSON
) ENGINE = ReplacingMergeTree()
      ORDER BY (`hash`);

CREATE MATERIALIZED VIEW IF NOT EXISTS transaction_consumer TO spacebox.transaction
AS
SELECT signatures,
       hash,
       height,
       success,
       messages,
       memo,
       signer_infos,
       fee,
       signer,
       gas_wanted,
       gas_used,
       raw_log,
       logs
FROM spacebox.transaction_topic
GROUP BY signatures, hash, height, success, messages, memo, signer_infos, fee, signer, gas_wanted, gas_used, raw_log,
         logs;
