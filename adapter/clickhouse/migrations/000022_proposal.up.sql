-- 000022_proposal.up.sql
CREATE TABLE IF NOT EXISTS spacebox.proposal
(
    `id`                Int64,
    `title`             String,
    `description`       String,
    `content`           json,
    `proposal_route`    String,
    `proposal_type`     String,
    `submit_time`       TIMESTAMP,
    `deposit_end_time`  TIMESTAMP,
    `voting_start_time` TIMESTAMP,
    `voting_end_time`   TIMESTAMP,
    `proposer_address`  String,
    `status`            String
) ENGINE = ReplacingMergeTree()
      ORDER BY (`id`);
