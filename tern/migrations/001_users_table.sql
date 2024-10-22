-- Write your migrate up statements here

CREATE TYPE role AS ENUM ('admin', 'editor', 'user');

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(24) UNIQUE NOT NULL,
    email VARCHAR(128) UNIQUE NOT NULL,
    pass_hash VARCHAR(60) NOT NULL,
    role role NOT NULL DEFAULT 'user',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

---- create above / drop below ----

DROP TABLE IF EXISTS users;
DROP TYPE IF EXISTS role;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
