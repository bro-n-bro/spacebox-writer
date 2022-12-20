--
CREATE TABLE my_test.block
(
    `height`           UInt256,
    `hash`             String,
    `num_txs`          Int256,
    `total_gas`        UInt256,
    `proposer_address` String,
    `timestamp`        TIMESTAMP
) ENGINE = MergeTree()
      PRIMARY KEY (height);


--
CREATE TABLE my_test.message
(
    `transaction_hash` String,
    `index`            UInt256,
    `type`             String,
    `value`            json
) ENGINE = MergeTree()
      PRIMARY KEY (transaction_hash);

--
CREATE TABLE my_test.transaction
(
    `hash`         String,
    `height`       UInt256,
    `success`      BOOL,
    `messages`     json,
    `memo`         String,
    `signatures`   String, -- TODO: // @malekvictor
    `signer_infos` json,
    `fee`          json,
    `signer`       String,
    `gas_wanted`   UInt256,
    `gas_used`     UInt256,
    `raw_log`      String,
    `logs`         json
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.account
(
    `address` String,
    `height`  UInt256
) ENGINE = MergeTree()
      PRIMARY KEY (height);


--
CREATE TABLE my_test.supply
(
    `height` UInt256,
    `coins`  json
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.account_balance
(
    `address` String,
    `height`  UInt256,
    `coins`   json
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.send_message
(
    `height`       UInt256,
    `address_from` String,
    `address_to`   String,
    `tx_hash`      String,
    `coins`        json
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.multisend_message
(
    `height`       UInt256,
    `address_from` String,
    `address_to`   String,
    `tx_hash`      String,
    `coins`        json
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.distribution_params
(
    `height` UInt256,
    `params` json
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.community_pool
(
    `height` UInt256,
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
    `height`           UInt256
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.delegation_reward
(
    operator_address  String,
    delegator_address String,
    withdraw_address  String,
    coins             json,
    height            UInt256
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.delegation_reward_message
(
    validator_address String,
    delegator_address String,
    coins             json,
    height            UInt256,
    tx_hash           String
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.gov_params
(
    deposit_params json,
    voting_params  json,
    tally_params   json,
    height         UInt256
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.proposal
(
    `id`                UInt256,
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
    proposal_id       UInt256,
    depositor_address String,
    coins             json,
    height            UInt256,
    tx_hash           String
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.proposal_tally_result
(
    proposal_id  UInt256,
    yes          UInt256,
    abstain      UInt256,
    no           UInt256,
    no_with_veto UInt256,
    height       UInt256
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.proposal_vote_message
(
    proposal_id   UInt256,
    voter_address String,
    option        UInt256,
    height        UInt256,
    tx_hash       String
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.mint_params
(
    params json,
    height UInt256
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.inflation
(
    height           UInt256,
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
    status            UInt256,
    jailed            BOOL,
    height            UInt256
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.validator_info
(
    consensus_address     String,
    operator_address      String,
    self_delegate_address String,
    min_self_delegation   UInt256,
    height                UInt256
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.staking_params
(
    `height` UInt256,
    `params` json
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.staking_pool
(
    `height`            UInt256,
    `not_bonded_tokens` UInt256,
    `bonded_tokens`     UInt256
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.redelegation
(
    `delegator_address`     String,
    `src_validator_address` String,
    `dst_validator_address` String,
    `coins`                 json,
    `height`                UInt256,
    `completion_time`       TIMESTAMP
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.redelegation_message
(
    `delegator_address`     String,
    `src_validator_address` String,
    `dst_validator_address` String,
    `coins`                 json,
    `height`                UInt256,
    `completion_time`       TIMESTAMP,
    `tx_hash`               String
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.unbonding_delegation
(
    `validator_address`    String,
    `delegator_address`    String,
    `coins`                json,
    `completion_timestamp` TIMESTAMP,
    `height`               UInt256
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.unbonding_delegation_message
(
    `validator_address`    String,
    `delegator_address`    String,
    `coin`                 json,
    `completion_timestamp` TIMESTAMP,
    `height`               UInt256,
    `tx_hash`              String
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.delegation
(
    `operator_address`  String,
    `delegator_address` String,
    `coin`              json,
    `height`            UInt256
) ENGINE = MergeTree()
      PRIMARY KEY (height);

--
CREATE TABLE my_test.delegation_message
(
    `operator_address`  String,
    `delegator_address` String,
    `coin`              json,
    `height`            UInt256,
    `tx_hash`           String
) ENGINE = MergeTree()
      PRIMARY KEY (height);