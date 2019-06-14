CREATE DATABASE mychannel;

use mychannel;

-- ----------------------------
--  Table structure for `transactions`
-- ----------------------------
DROP TABLE IF EXISTS transactions;

CREATE TABLE transactions (
  txid character varying(256) PRIMARY KEY,
  type integer DEFAULT NULL,
  date integer DEFAULT NULL,
  time Timestamp DEFAULT NULL,
  sender character varying(64) DEFAULT NULL,
  receiver character varying(64) DEFAULT NULL,
  amount decimal(64,2) DEFAULT NULL,
  token character varying(64) DEFAULT NULL
);