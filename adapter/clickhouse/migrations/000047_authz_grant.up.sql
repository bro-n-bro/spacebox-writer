-- 000047_authz_grant.up.sql
CREATE TABLE IF NOT EXISTS spacebox.authz_grant
(
    `height`          Int64,
    `granter_address` String,
    `grantee_address` String,
    `msg_type`        String,
    `expiration`      TIMESTAMP

) ENGINE = ReplacingMergeTree()
      ORDER BY (`grantee_address`, `granter_address`, `msg_type`);