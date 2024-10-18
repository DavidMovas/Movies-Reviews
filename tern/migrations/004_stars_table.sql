-- Write your migrate up statements here

CREATE TABLE stars (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    middle_name VARCHAR(50),
    last_name VARCHAR(50) NOT NULL,
    birth_date DATE NOT NULL,
    birth_place VARCHAR(100),
    death_date DATE,
    bio TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

---- create above / drop below ----

DROP TABLE stars;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
