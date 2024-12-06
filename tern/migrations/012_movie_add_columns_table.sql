-- Write your migrate up statements here

ALTER TABLE movies
ADD COLUMN poster_url VARCHAR(255),
ADD COLUMN imdb_rating float4,
ADD COLUMN imdb_url VARCHAR(255),
ADD COLUMN  metascore INT,
ADD COLUMN metascore_url VARCHAR(255);

---- create above / drop below ----

ALTER TABLE movies
DROP COLUMN poster_url,
DROP COLUMN imdb_rating,
DROP COLUMN imdb_url,
DROP COLUMN metascore,
DROP COLUMN metascore_url;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
