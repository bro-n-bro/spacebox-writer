-- 000015_supply.up.sql
CREATE TABLE IF NOT EXISTS spacebox.supply
(
    `coins`  json,
    `height` Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY ( `height`);
