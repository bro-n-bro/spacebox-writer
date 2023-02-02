-- 000002_account_balance.up.sql
CREATE TABLE IF NOT EXISTS spacebox.account_balance
(
    `address` String,
    `height`  Int64,
    `coins`   json
) ENGINE = ReplacingMergeTree(`height`)
      ORDER BY (`address`);
