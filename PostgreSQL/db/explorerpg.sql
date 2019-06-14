--
--    SPDX-License-Identifier: Apache-2.0
--

CREATE USER :user WITH PASSWORD ':passwd';
DROP DATABASE IF EXISTS :dbname;
CREATE DATABASE :dbname owner :user;
\c :dbname;
--


-- ----------------------------
--  Table structure for `transactions`
-- ----------------------------
DROP TABLE IF EXISTS transactions;

CREATE TABLE transactions (
  -- id SERIAL PRIMARY KEY,
  txid character varying(256) PRIMARY KEY,
  type integer DEFAULT NULL,
  time Timestamp DEFAULT NULL,
  sender character varying(64) DEFAULT NULL,
  receiver character varying(64) DEFAULT NULL,
  amount decimal(64,2) DEFAULT NULL,
  token character varying(64) DEFAULT NULL
);

ALTER table transactions owner to :user;

DROP INDEX IF EXISTS transactions_txid_idx;
CREATE INDEX ON Transactions (txid);

DROP INDEX IF EXISTS transactions_receiver_idx;
CREATE INDEX ON Transactions (receiver);

GRANT SELECT, INSERT, UPDATE,DELETE ON ALL TABLES IN SCHEMA PUBLIC to :user;
