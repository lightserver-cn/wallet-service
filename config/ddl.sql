DROP TABLE IF EXISTS "t_transaction";
DROP SEQUENCE IF EXISTS transaction_id_seq;
CREATE SEQUENCE transaction_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."t_transaction"
(
    "id"                 integer        DEFAULT nextval('transaction_id_seq') NOT NULL,
    "sender_wallet_id"   integer        DEFAULT '0',
    "receiver_wallet_id" integer        DEFAULT '0',
    "amount"             numeric(15, 2) DEFAULT '0.00'                        NOT NULL,
    "transaction_type"   smallint       DEFAULT '0'                           NOT NULL,
    "created_at"         timestamp      DEFAULT CURRENT_TIMESTAMP             NOT NULL,
    CONSTRAINT "transaction_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

CREATE INDEX "transaction_receiver_wallet_id" ON "public"."t_transaction" USING btree ("receiver_wallet_id");

CREATE INDEX "transaction_sender_wallet_id" ON "public"."t_transaction" USING btree ("sender_wallet_id");

COMMENT
ON COLUMN "public"."t_transaction"."transaction_type" IS '1-deposit, 2-withdraw, 3-transfer';


DROP TABLE IF EXISTS "t_user";
DROP SEQUENCE IF EXISTS user_id_seq;
CREATE SEQUENCE user_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."t_user"
(
    "id"            integer   DEFAULT nextval('user_id_seq') NOT NULL,
    "username"      character varying(255)                   NOT NULL,
    "email"         character varying(255)                   NOT NULL,
    "password_hash" character varying(255)                   NOT NULL,
    "status"        smallint  DEFAULT '1'                    NOT NULL,
    "created_at"    timestamp DEFAULT CURRENT_TIMESTAMP      NOT NULL,
    "updated_at"    timestamp DEFAULT CURRENT_TIMESTAMP      NOT NULL,
    CONSTRAINT "user_email" UNIQUE ("email"),
    CONSTRAINT "user_pkey" PRIMARY KEY ("id"),
    CONSTRAINT "user_username" UNIQUE ("username")
) WITH (oids = false);

COMMENT
ON COLUMN "public"."t_user"."status" IS '1-Valid, 2-Invalid, 3-Disabled';


DROP TABLE IF EXISTS "t_wallet";
DROP SEQUENCE IF EXISTS wallet_id_seq;
CREATE SEQUENCE wallet_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "public"."t_wallet"
(
    "id"         integer        DEFAULT nextval('wallet_id_seq') NOT NULL,
    "uid"        integer        DEFAULT '0'                      NOT NULL,
    "balance"    numeric(15, 2) DEFAULT '0.00'                   NOT NULL,
    "created_at" timestamp      DEFAULT CURRENT_TIMESTAMP        NOT NULL,
    "updated_at" timestamp      DEFAULT CURRENT_TIMESTAMP        NOT NULL,
    CONSTRAINT "wallet_pkey" PRIMARY KEY ("id")
) WITH (oids = false);

CREATE INDEX "wallet_uid" ON "public"."t_wallet" USING btree ("uid");

DROP DATABASE IF EXISTS test_postgres;
CREATE DATABASE test_postgres;
DROP DATABASE IF EXISTS db_test;
CREATE DATABASE db_test;
