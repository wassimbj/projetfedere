-- Adminer 4.8.1 PostgreSQL 13.4 (Debian 13.4-4.pgdg110+1) dump

DROP TABLE IF EXISTS "blocked_users";
DROP SEQUENCE IF EXISTS bu_id_seq;
CREATE SEQUENCE bu_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."blocked_users" (
    "id" integer DEFAULT nextval('bu_id_seq') NOT NULL,
    "user_blocked" integer NOT NULL,
    "blocked_by" integer NOT NULL,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT "bu_pkey" PRIMARY KEY ("id")
) WITH (oids = false);


DROP TABLE IF EXISTS "messages";
DROP SEQUENCE IF EXISTS msg_id_seq;
CREATE SEQUENCE msg_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."messages" (
    "id" integer DEFAULT nextval('msg_id_seq') NOT NULL,
    "sent_from" integer NOT NULL,
    "sent_to" integer NOT NULL,
    "msg" character varying(1000),
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT "msgs_pkey" PRIMARY KEY ("id")
) WITH (oids = false);


DROP TABLE IF EXISTS "users";
DROP SEQUENCE IF EXISTS users_id_seq;
CREATE SEQUENCE users_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."users" (
    "id" integer DEFAULT nextval('users_id_seq') NOT NULL,
    "firstname" character varying(50) NOT NULL,
    "lastname" character varying(50) NOT NULL,
    "username" character varying(50),
    "email" character varying(20) NOT NULL,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "is_email_confirmed" boolean DEFAULT false NOT NULL,
    "password" character varying(70) NOT NULL,
    CONSTRAINT "usernames" UNIQUE ("username"),
    CONSTRAINT "users_email" UNIQUE ("email"),
    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
) WITH (oids = false);


-- 2022-04-24 07:09:10.384747+00
