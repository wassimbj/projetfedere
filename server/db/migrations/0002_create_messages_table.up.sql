CREATE SEQUENCE IF NOT EXISTS msg_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE IF NOT EXISTS "messages" (
    "id" integer DEFAULT nextval('msg_id_seq') NOT NULL,
    "sent_from" integer NOT NULL,
    "sent_to" integer NOT NULL,
    "msg" varchar(1000),
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT "msgs_pkey" PRIMARY KEY ("id")
) WITH (oids = false);
