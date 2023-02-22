-- 000021_gov_params.up.sql
CREATE TABLE IF NOT EXISTS spacebox.gov_params_topic
(
    `deposit_params` String, -- TODO: @malekviktor this is json
    `voting_params`  String, -- TODO: @malekviktor this is json
    `tally_params`   String, -- TODO: @malekviktor this is json
    `height`         Int64
) ENGINE = Kafka('kafka:9093', 'community_pool', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.gov_params
(
    `deposit_params` String, -- TODO: @malekviktor this is json
    `voting_params`  String, -- TODO: @malekviktor this is json
    `tally_params`   String, -- TODO: @malekviktor this is json
    `height`         Int64
) ENGINE = ReplacingMergeTree()
      ORDER BY (`height`);

CREATE MATERIALIZED VIEW IF NOT EXISTS gov_params_consumer TO spacebox.gov_params
AS
SELECT deposit_params, voting_params, tally_params, height
FROM spacebox.gov_params_topic
GROUP BY deposit_params, voting_params, tally_params, height;
