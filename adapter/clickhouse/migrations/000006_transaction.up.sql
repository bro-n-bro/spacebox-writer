-- 000006_transaction.up.sql
CREATE TABLE IF NOT EXISTS spacebox.transaction
(
    `signatures` Array(String),
    `hash`         String,
    `height`       Int64,
    `success`      BOOL,
    `messages`     String,
    `memo`         String,
    `signer_infos` json,
    `fee`          json,
    `signer`       String,
    `gas_wanted`   Int64,
    `gas_used`     Int64,
    `raw_log`      String,
    `logs`         json
) ENGINE = ReplacingMergeTree()
      ORDER BY (`hash`);