-- 000038_create_validator_message.up.sql
CREATE TABLE IF NOT EXISTS spacebox.create_validator_message
(
    `height`              Int64,
    `msg_index`           Int64,
    `tx_hash`             String,
    `delegator_address`   String,
    `validator_address`   String,
    `description`         String,
    `commission_rates`    Float64,
    `min_self_delegation` Int64,
    `pubkey`              String
) ENGINE = ReplacingMergeTree()
      ORDER BY (`tx_hash`, `msg_index`);
