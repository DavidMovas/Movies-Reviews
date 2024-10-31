CREATE TABLE reviews (
    id SERIAL PRIMARY KEY,
    movie_id INTEGER NOT NULL REFERENCES movies(id),
    user_id INTEGER NOT NULL REFERENCES users(id),
    rating SMALLINT,
    title VARCHAR(100),
    content VARCHAR(1500),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,

    UNIQUE (movie_id, user_id)
);
CREATE INDEX idx_reviews_movie_id ON reviews (movie_id);
CREATE INDEX idx_reviews_user_id ON reviews (user_id);
---- create above / drop below ----
DROP INDEX idx_reviews_user_id;
DROP INDEX idx_reviews_movie_id;
DROP TABLE reviews;
