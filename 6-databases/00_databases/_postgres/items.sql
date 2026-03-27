SET TIME ZONE '+00:00';

DROP TABLE IF EXISTS items;

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    updated VARCHAR(255)
);

INSERT INTO items (id, title, description, updated) VALUES
    (1,	'database/sql',	'Рассказать про базы данных',	'rvasily'),
    (2,	'memcache',	'Рассказать про мемкеш с примером использования',	NULL);

SELECT setval(pg_get_serial_sequence('items', 'id'), MAX(id)) FROM items;
