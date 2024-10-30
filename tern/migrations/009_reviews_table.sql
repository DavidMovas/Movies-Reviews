CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    movie_id INTEGER NOT NULL REFERENCES movies(id) UNIQUE,
    user_id INTEGER NOT NULL REFERENCES users(id) UNIQUE,
    rating SMALLINT,
    title VARCHAR(100),
    content VARCHAR(1500),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_ar TIMESTAMP,
    deleted_at TIMESTAMP
);
CREATE INDEX idx_reviews_movie_id ON reviews (movie_id);
CREATE INDEX idx_reviews_user_id ON reviews (user_id);
---- create above / drop below ----
DROP INDEX idx_reviews_user_id;
DROP INDEX idx_reviews_movie_id;
DROP TABLE reviews;
