CREATE SEQUENCE IF NOT EXISTS bu_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE IF NOT EXISTS "blocked_users" (
   "id" integer DEFAULT nextval('bu_id_seq') NOT NULL,
   "user_blocked" integer NOT NULL,
   "blocked_by" integer NOT NULL,
   "created_at" timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
   CONSTRAINT "bu_pkey" PRIMARY KEY ("id")
) WITH (oids = false);
