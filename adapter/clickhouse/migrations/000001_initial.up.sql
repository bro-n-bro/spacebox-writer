-- ## Modules

-- ### Auth

-- #### Account
-- ```sql=
CREATE TABLE spacebox.account
(
    `address` String,
    `height`  Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);
-- ```

-- ### Bank

-- #### Supply
-- ```sql=
CREATE TABLE spacebox.supply
(
    `height` Int64,
    `coins`  json
) ENGINE = MergeTree()
      PRIMARY KEY (height);
-- ```

-- #### Account Balance
-- ```sql=
CREATE TABLE spacebox.account_balance
(
    `address` String,
    `height`  Int64,
    `coins`   json
) ENGINE = MergeTree()
      PRIMARY KEY (height);
-- ```

-- #### Send Message
-- ```sql=
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
-- ```

-- #### Multisend Message
-- ```sql=
CREATE TABLE spacebox.multisend_message
(
    `height`       Int64,
    `address_from` String,
    `addresses_to` json,
    `tx_hash`      String,
    `coins`        json,
    `msg_index`    Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);
-- ```

-- ### Core

-- #### Block
-- ```sql=
CREATE TABLE spacebox.block
(
    `height`           Int64,
    `hash`             String,
    `num_txs`          Int64,
    `total_gas`        Int64,
    `proposer_address` String,
    `timestamp`        TIMESTAMP
) ENGINE = MergeTree()
      PRIMARY KEY (height);
-- ```

-- #### Message
-- ```sql=
CREATE TABLE spacebox.message
(
    `transaction_hash` String,
    `index`            Int64,
    `type`             String,
    `value`            json
) ENGINE = MergeTree()
      PRIMARY KEY (transaction_hash);
-- ```

-- #### Transaction
-- ```sql=
CREATE TABLE spacebox.transaction
(
    `hash`         String,
    `height`       Int64,
    `success`      BOOL,
    `messages`     json,
    `memo`         String,
    `signatures`   json,
    `signer_infos` json,
    `fee`          json,
    `signer`       String,
    `gas_wanted`   Int64,
    `gas_used`     Int64,
    `raw_log`      String,
    `logs`         json,
    `code`         UInt32
) ENGINE = MergeTree()
      PRIMARY KEY (height);
-- ```

-- ### Distribution

-- #### Distribution Params
-- ```sql=
CREATE TABLE spacebox.distribution_params
(
    `height` Int64,
    `params` json
) ENGINE = MergeTree()
      PRIMARY KEY (height);
-- ```

-- #### Community Pool
-- ```sql=
CREATE TABLE spacebox.community_pool
(
    `height` Int64,
    `coins`  json
) ENGINE = MergeTree()
      PRIMARY KEY (height);
-- ```

-- #### Validator Commission
-- ```sql=
CREATE TABLE spacebox.validator_commission
(
    `operator_address` String,
    `commission`       Float64,
    `max_change_rate`  Float64,
    `max_rate`         Float64,
    `height`           Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);
-- ```

-- #### Delegation Reward Message
-- ```sql=
CREATE TABLE spacebox.delegation_reward_message
(
    `validator_address` String,
    `delegator_address` String,
    `coins`             json,
    `height`            Int64,
    `tx_hash`           String,
    `msg_index`       Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);
-- ```

-- ### Gov

-- #### Gov Params
-- ```sql=
CREATE TABLE spacebox.gov_params
(
    `deposit_params` json,
    `voting_params`  json,
    `tally_params`   json,
    `height`         Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);
-- ```

-- #### Proposal
-- ```sql=
CREATE TABLE spacebox.proposal
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
) ENGINE = MergeTree()
      PRIMARY KEY (id);
-- ```

-- #### Proposal Deposit Message
-- ```sql=
CREATE TABLE spacebox.proposal_deposit_message
(
    `proposal_id`       Int64,
    `depositor_address` String,
    `coins`             json,
    `height`            Int64,
    `tx_hash`           String,
    `msg_index`       Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);
-- ```


-- #### Proposal Tally Result
-- ```sql=
CREATE TABLE spacebox.proposal_tally_result
(
    `proposal_id`  Int64,
    `yes`          Int64,
    `abstain`      Int64,
    `no`           Int64,
    `no_with_veto` Int64,
    `height`       Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);
-- ```

-- #### Proposal Vote Message
-- ```sql=
CREATE TABLE spacebox.proposal_vote_message
(
    `proposal_id`   Int64,
    `voter_address` String,
    `option`        String,
    `height`        Int64,
    `tx_hash`       String,
    `msg_index`   Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);
-- ```

-- ### Mint

-- #### Mint Params
-- ```sql=
CREATE TABLE spacebox.mint_params
(
    `params` json,
    `height` Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);
-- ```

-- #### Annual Provision
-- ```sql=
CREATE TABLE spacebox.annual_provision
(
    `height`           Int64,
    `annual_provision` Float64
) ENGINE = MergeTree()
      PRIMARY KEY (height);
-- ```

-- ### Staking

-- #### Validator
-- ```sql=
CREATE TABLE spacebox.validator
(
    `consensus_address` String,
    `operator_address`  String,
    `consensus_pubkey`  String
) ENGINE = MergeTree()
      PRIMARY KEY (consensus_address);
-- ```

-- #### Validator Status
-- ```sql=
CREATE TABLE spacebox.validator_status
(
    `validator_address` String,
    `status`            Int64,
    `jailed`            BOOL,
    `height`            Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);
-- ```

-- #### Validator Info
-- ```sql=
CREATE TABLE spacebox.validator_info
(
    `consensus_address`     String,
    `operator_address`      String,
    `self_delegate_address` String,
    `min_self_delegation`   Int64,
    `height`                Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);
-- ```

-- #### Staking Params
-- ```sql=
CREATE TABLE spacebox.staking_params
(
    `height` Int64,
    `params` json
) ENGINE = MergeTree()
      PRIMARY KEY (height);
-- ```

-- #### Staking Pool
-- ```sql=
CREATE TABLE spacebox.staking_pool
(
    `height`            Int64,
    `not_bonded_tokens` Int64,
    `bonded_tokens`     Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);
-- ```

-- #### Redelegation
-- ```sql=
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
-- ```

-- #### Redelegation Message
-- ```sql=
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
-- ```

-- #### Unbonding Delegation
-- ```sql=
CREATE TABLE spacebox.unbonding_delegation
(
    `validator_address`    String,
    `delegator_address`    String,
    `coin`                 json,
    `completion_timestamp` TIMESTAMP,
    `height`               Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);
-- ```

-- #### Unbonding Delegation Message
-- ```sql=
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
-- ```

-- #### Delegation
-- ```sql=
CREATE TABLE spacebox.delegation
(
    `operator_address`  String,
    `delegator_address` String,
    `coin`              json,
    `height`            Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);
-- ```

-- #### Delegation Message
-- ```sql=
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
-- ```

-- #### Validator Description
-- ```sql=
CREATE TABLE spacebox.validator_description
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
-- ```