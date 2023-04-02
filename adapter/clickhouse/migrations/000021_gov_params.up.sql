-- 000021_gov_params.up.sql
CREATE TABLE IF NOT EXISTS spacebox.gov_params
(
    `deposit_params` json,
    `voting_params`  json,
    `tally_params`   json,
    `height`         Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`height`);