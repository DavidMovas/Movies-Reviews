-- Write your migrate up statements here

ALTER TABLE stars
ADD COLUMN avatar_url VARCHAR(255),
ADD COLUMN imdb_url VARCHAR(255);

---- create above / drop below ----

ALTER TABLE stars
DROP COLUMN avatar_url,
DROP COLUMN imdb_url;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
