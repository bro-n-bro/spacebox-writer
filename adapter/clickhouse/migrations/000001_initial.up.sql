CREATE TABLE spacebox.account
(
    `address` String,
    `height`  Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height); -- TODO

-- 0001_account.sql
-- 0002_account_balance.sql
-- 0003_block.sql
-- 0004_community_pool.sql
-- 0005_delegation_reward_message.sql
-- 0006_distribution_params.sql
-- 0007_gov_params.sql
-- 0008_message.sql
-- 0009_multisend_message.sql
-- 0010_proposal.sql
-- 0011_proposal_deposit_message.sql
-- 0012_proposal_tally_result.sql
-- 0013_proposal_vote_message.sql
-- 0014_send_message.sql
-- 0015_supply.sql
-- 0016_transaction.sql
-- 0017_validator_commission.sql
-- 0018_validator.sql
-- 0019_validator_delegation.sql
-- 0020_validator_unbonding_delegation.sql
-- 0021_validator_unbonding_delegation_message

CREATE TABLE spacebox.supply
(
    `height` Int64,
    `coins`  json
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.account_balance
(
    `address` String,
    `height`  Int64,
    `coins`   json
) ENGINE = MergeTree()
      PRIMARY KEY (address, height);

CREATE TABLE spacebox.send_message
(
    `height`       Int64,
    `address_from` String,
    `address_to`   String,
    `tx_hash`      String,
    `coins`        json,
    `msg_index`    Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.multisend_message
(
    `height`       Int64,
    `address_from` String,
    `addresses_to` Array(String),
    `tx_hash`      String,
    `coins`        json,
    `msg_index`    Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.block -- TODO
(
    `height`           Int64,
    `hash`             String,
    `num_txs`          Int64,
    `total_gas`        Int64,
    `proposer_address` String,
    `timestamp`        TIMESTAMP
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.message -- TODO
(
    `transaction_hash` String,
    `msg_index`        Int64,
    `type`             String,
    `signer`           String,
    `value`            String,
    `involved_accounts_addresses` Array(String)
) ENGINE = MergeTree()
      PRIMARY KEY (transaction_hash);

CREATE TABLE spacebox.transaction
(
    `hash`         String,
    `height`       Int64,
    `success`      BOOL,
    `messages`     String,
    `memo`         String,
    `signatures`   Array(String),
    `signer_infos` json,
    `fee`          json,
    `signer`       String,
    `gas_wanted`   Int64,
    `gas_used`     Int64,
    `raw_log`      String,
    `logs`         json
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.distribution_params
(
    `height` Int64,
    `params` json
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.community_pool
(
    `height` Int64,
    `coins`  json
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.validator_commission -- TODO
(
    `operator_address` String,
    `commission`       Float64,
    `max_change_rate`  Float64,
    `max_rate`         Float64,
    `height`           Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.delegation_reward_message
(
    `validator_address` String,
    `delegator_address` String,
    `coins`             json,
    `height`            Int64,
    `tx_hash`           String,
    `msg_index`         Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.gov_params
(
    `deposit_params` json,
    `voting_params`  json,
    `tally_params`   json,
    `height`         Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.proposal
(
    `id`                UInt64,
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
) ENGINE = MergeTree()
      PRIMARY KEY (id);

CREATE TABLE spacebox.proposal_deposit_message
(
    `proposal_id`       UInt64,
    `depositor_address` String,
    `coins`             json,
    `height`            Int64,
    `tx_hash`           String,
    `msg_index`         Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.proposal_tally_result -- TODO
(
    `proposal_id`  UInt64,
    `yes`          Int64,
    `abstain`      Int64,
    `no`           Int64,
    `no_with_veto` Int64,
    `height`       Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.proposal_vote_message -- TODO
(
    `proposal_id`   UInt64,
    `voter_address` String,
    `option`        String,
    `height`        Int64,
    `tx_hash`       String,
    `msg_index`     Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.mint_params
(
    `params` json,
    `height` Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.annual_provision -- TODO
(
    `height`           Int64,
    `annual_provision` Float64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.validator -- TODO
(
    `consensus_address` String,
    `operator_address`  String,
    `consensus_pubkey`  String
) ENGINE = MergeTree()
      PRIMARY KEY (consensus_address);

CREATE TABLE spacebox.validator_status -- TODO
(
    `validator_address` String,
    `status`            Int64,
    `jailed`            BOOL,
    `height`            Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.validator_info -- TODO
(
    `consensus_address`     String,
    `operator_address`      String,
    `self_delegate_address` String,
    `min_self_delegation`   Int64,
    `height`                Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.staking_params
(
    `height` Int64,
    `params` json
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.staking_pool -- TODO
(
    `height`            Int64,
    `not_bonded_tokens` Int64,
    `bonded_tokens`     Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.redelegation
(
    `delegator_address`     String,
    `src_validator_address` String,
    `dst_validator_address` String,
    `coin`                  json,
    `height`                Int64,
    `completion_time`       TIMESTAMP
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.redelegation_message
(
    `delegator_address`     String,
    `src_validator_address` String,
    `dst_validator_address` String,
    `coin`                  json,
    `height`                Int64,
    `completion_time`       TIMESTAMP,
    `tx_hash`               String,
    `msg_index`             Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.unbonding_delegation
(
    `validator_address`    String,
    `delegator_address`    String,
    `coin`                 json,
    `completion_timestamp` TIMESTAMP,
    `height`               Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.unbonding_delegation_message
(
    `validator_address`    String,
    `delegator_address`    String,
    `coin`                 json,
    `completion_timestamp` TIMESTAMP,
    `height`               Int64,
    `tx_hash`              String,
    `msg_index`            Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.delegation
(
    `operator_address`  String,
    `delegator_address` String,
    `coin`              json,
    `height`            Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.delegation_message
(
    `operator_address`  String,
    `delegator_address` String,
    `coin`              json,
    `height`            Int64,
    `tx_hash`           String,
    `msg_index`         Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

CREATE TABLE spacebox.validator_description -- TODO
(
    `validator_address` String,
    `moniker`           String,
    `identity`          String,
    `avatar_url`        String,
    `website`           String,
    `security_contact`  String,
    `details`           String,
    `height`            Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);