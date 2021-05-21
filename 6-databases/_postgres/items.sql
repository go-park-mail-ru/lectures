DROP TABLE IF EXISTS "items";
DROP SEQUENCE IF EXISTS items_id_seq;
CREATE SEQUENCE items_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 START 1 CACHE 1;

CREATE TABLE "items" (
    "id" integer DEFAULT nextval('items_id_seq') NOT NULL,
    "title" character varying(255) NOT NULL,
    "description" text NOT NULL,
    "updated" character varying(255)
) WITH (oids = false);

INSERT INTO items (id, title, description, updated) VALUES
(1,	'database/sql',	'Рассказать про базы данных',	'rvasily'),
(2,	'memcache',	'Рассказать про мемкеш с примером использования',	NULL);
