-- 000018_validator_description.up.sql
CREATE TABLE IF NOT EXISTS spacebox.validator_description_topic
(
    `validator_address` String,
    `moniker`           String,
    `identity`          String,
    `avatar_url`        String,
    `website`           String,
    `security_contact`  String,
    `details`           String,
    `height`            Int64
) ENGINE = Kafka('kafka:9093', 'validator_description', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.validator_description
(
    `validator_address` String,
    `moniker`           String,
    `identity`          String,
    `avatar_url`        String,
    `website`           String,
    `security_contact`  String,
    `details`           String,
    `height`            Int64
) ENGINE = ReplacingMergeTree(`height`)
      ORDER BY (`validator_address`);

CREATE MATERIALIZED VIEW IF NOT EXISTS validator_description_consumer TO spacebox.validator_description
AS
SELECT validator_address, moniker, identity, avatar_url, website, security_contact, details, height
FROM spacebox.validator_description_topic
GROUP BY validator_address, moniker, identity, avatar_url, website, security_contact, details, height;
