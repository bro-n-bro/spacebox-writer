-- 000022_proposal.up.sql
CREATE TABLE IF NOT EXISTS spacebox.proposal_topic
(
    `id`                Int64,
    `title`             String,
    `description`       String,
    `content`           String, -- TODO: @malekviktor this is json
    `proposal_route`    String,
    `proposal_type`     String,
    `submit_time`       TIMESTAMP,
    `deposit_end_time`  TIMESTAMP,
    `voting_start_time` TIMESTAMP,
    `voting_end_time`   TIMESTAMP,
    `proposer_address`  String,
    `status`            String
) ENGINE = Kafka('kafka:9093', 'proposal', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.proposal
(
    `id`                Int64,
    `title`             String,
    `description`       String,
    `content`           String, -- TODO: @malekviktor this is json
    `proposal_route`    String,
    `proposal_type`     String,
    `submit_time`       TIMESTAMP,
    `deposit_end_time`  TIMESTAMP,
    `voting_start_time` TIMESTAMP,
    `voting_end_time`   TIMESTAMP,
    `proposer_address`  String,
    `status`            String
) ENGINE = ReplacingMergeTree()
      ORDER BY (`id`, `status`);

CREATE MATERIALIZED VIEW IF NOT EXISTS proposal_consumer TO spacebox.proposal
AS
SELECT id,
       title,
       description,
       content,
       proposal_route,
       proposal_type,
       submit_time,
       deposit_end_time,
       voting_start_time,
       voting_end_time,
       proposer_address,
       status
FROM spacebox.proposal_topic
GROUP BY id, title, description, content, proposal_route, proposal_type, submit_time, deposit_end_time,
         voting_start_time, voting_end_time, proposer_address, status;
