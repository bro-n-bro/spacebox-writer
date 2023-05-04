-- 000032_withdraw_validator_commission_message.up.sql
CREATE TABLE IF NOT EXISTS spacebox.withdraw_validator_commission_message
(
    `height`              Int64,
    `tx_hash`             String,
    `msg_index`           Int64,
    `withdraw_commission` json,
    `validator_address`   String
) ENGINE = ReplacingMergeTree()
      ORDER BY (`tx_hash`, `msg_index`);
