-- 000031_delegation_message.up.sql
CREATE TABLE IF NOT EXISTS spacebox.delegation_message
(
    `operator_address`  String,
    `delegator_address` String,
    `coin`              json,
    `height`            Int64,
    `tx_hash`           String,
    `msg_index`         Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`tx_hash`, `msg_index`);
