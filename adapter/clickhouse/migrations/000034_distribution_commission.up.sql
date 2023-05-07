-- 000034_distribution_commission.up.sql
CREATE TABLE IF NOT EXISTS spacebox.distribution_commission
(
    `validator` String,
    `amount`    json,
    `height`    Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`height`, `validator`);
