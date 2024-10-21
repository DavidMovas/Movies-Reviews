-- Write your migrate up statements here

CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    description TEXT NOT NULL,
    release_date DATE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

---- create above / drop below ----

DROP TABLE movies;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
