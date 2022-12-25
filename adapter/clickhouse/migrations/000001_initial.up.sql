--
CREATE TABLE my_test.block
(
    `height`           Int64,
    `hash`             String,
    `num_txs`          Int64,
    `total_gas`        Int64,
    `proposer_address` String,
    `timestamp`        TIMESTAMP
) ENGINE = MergeTree()
      PRIMARY KEY (height);


--
CREATE TABLE my_test.message
(
    `transaction_hash` String,
    `index`            Int64,
    `type`             String,
    `value`            json
) ENGINE = MergeTree()
      PRIMARY KEY (transaction_hash);

--
CREATE TABLE my_test.transaction
(
    `hash`         String,
    `height`       Int64,
    `success`      BOOL,
    `messages`     json,
    `memo`         String,
    `signatures`   String, -- TODO: // @malekvictor
    `signer_infos` json,
    `fee`          json,
    `signer`       String,
    `gas_wanted`   Int64,
    `gas_used`     Int64,
    `raw_log`      String,
    `logs`         json
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.account
(
    `address` String,
    `height`  Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);


--
CREATE TABLE my_test.supply
(
    `height` Int64,
    `coins`  json
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.account_balance
(
    `address` String,
    `height`  Int64,
    `coins`   json
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.send_message
(
    `height`       Int64,
    `address_from` String,
    `address_to`   String,
    `tx_hash`      String,
    `coins`        json
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.multisend_message
(
    `height`       Int64,
    `address_from` String,
    `address_to`   String,
    `tx_hash`      String,
    `coins`        json
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.distribution_params
(
    `height` Int64,
    `params` json
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.community_pool
(
    `height` Int64,
    `coins`  json
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.validator_commission
(
    `operator_address` String,
    `commission`       Float64,
    `max_change_rate`  Float64,
    `max_rate`         Float64,
    `height`           Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.delegation_reward
(
    operator_address  String,
    delegator_address String,
    withdraw_address  String,
    coins             json,
    height            Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.delegation_reward_message
(
    validator_address String,
    delegator_address String,
    coins             json,
    height            Int64,
    tx_hash           String
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.gov_params
(
    deposit_params json,
    voting_params  json,
    tally_params   json,
    height         Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.proposal
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

--
CREATE TABLE my_test.proposal_deposit
(
    proposal_id       Int64,
    depositor_address String,
    coins             json,
    height            Int64,
    tx_hash           String
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.proposal_tally_result
(
    proposal_id  Int64,
    yes          Int64,
    abstain      Int64,
    no           Int64,
    no_with_veto Int64,
    height       Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.proposal_vote_message
(
    proposal_id   Int64,
    voter_address String,
    option        Int64,
    height        Int64,
    tx_hash       String
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.mint_params
(
    params json,
    height Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.inflation
(
    height           Int64,
    inflation        Float64,
    block_provision  json,
    annual_provision json
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.validator
(
    consensus_address String,
    consensus_pubkey  String
) ENGINE = MergeTree()
      PRIMARY KEY (consensus_address);

--
CREATE TABLE my_test.validator_status
(
    validator_address String,
    status            Int64,
    jailed            BOOL,
    height            Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.validator_info
(
    consensus_address     String,
    operator_address      String,
    self_delegate_address String,
    min_self_delegation   Int64,
    height                Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.staking_params
(
    `height` Int64,
    `params` json
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.staking_pool
(
    `height`            Int64,
    `not_bonded_tokens` Int64,
    `bonded_tokens`     Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.redelegation
(
    `delegator_address`     String,
    `src_validator_address` String,
    `dst_validator_address` String,
    `coin`                 json,
    `height`                Int64,
    `completion_time`       TIMESTAMP
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.redelegation_message
(
    `delegator_address`     String,
    `src_validator_address` String,
    `dst_validator_address` String,
    `coin`                 json,
    `height`                Int64,
    `completion_time`       TIMESTAMP,
    `tx_hash`               String
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.unbonding_delegation
(
    `validator_address`    String,
    `delegator_address`    String,
    `coin`                json,
    `completion_timestamp` TIMESTAMP,
    `height`               Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.unbonding_delegation_message
(
    `validator_address`    String,
    `delegator_address`    String,
    `coin`                 json,
    `completion_timestamp` TIMESTAMP,
    `height`               Int64,
    `tx_hash`              String
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.delegation
(
    `operator_address`  String,
    `delegator_address` String,
    `coin`              json,
    `height`            Int64
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.delegation_message
(
    `operator_address`  String,
    `delegator_address` String,
    `coin`              json,
    `height`            Int64,
    `tx_hash`           String
) ENGINE = MergeTree()
      PRIMARY KEY (height);