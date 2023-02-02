-- 000010_annual_provision.up.sql
CREATE TABLE IF NOT EXISTS spacebox.annual_provision_topic
(
    `height`           Int64,
    `annual_provision` Float64
) ENGINE = Kafka('kafka:9093', 'annual_provision', 'spacebox', 'JSONEachRow');

CREATE TABLE IF NOT EXISTS spacebox.annual_provision
(
    `height`           Int64,
    `annual_provision` Float64
) ENGINE = ReplacingMergeTree(`height`)
      ORDER BY (`height`);

CREATE MATERIALIZED VIEW IF NOT EXISTS annual_provision_consumer TO spacebox.annual_provision
AS
SELECT height, annual_provision
FROM spacebox.annual_provision_topic
GROUP BY height, annual_provision;
