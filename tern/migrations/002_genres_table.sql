-- Write your migrate up statements here

CREATE TABLE genres (
    id SERIAL PRIMARY KEY,
    name VARCHAR(32) NOT NULL
);

---- create above / drop below ----

DROP TABLE genres;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
