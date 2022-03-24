CREATE SEQUENCE IF NOT EXISTS users_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE IF NOT EXISTS "users" (
    "id" integer DEFAULT nextval('users_id_seq') NOT NULL,
    "firstname" character varying(50) NOT NULL,
    "lastname" character varying(50) NOT NULL,
    "username" character varying(50) NOT NULL,
    "email" character varying(20) NOT NULL,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    "is_email_confirmed" boolean DEFAULT 'false' NOT NULL,
    "password" character varying(70) NOT NULL,
    CONSTRAINT "users_email" UNIQUE ("email"),
    CONSTRAINT "usernames" UNIQUE ("username"),
    CONSTRAINT "users_pkey" PRIMARY KEY ("id")
) WITH (oids = false);
