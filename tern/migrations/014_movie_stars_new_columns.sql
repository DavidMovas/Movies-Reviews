-- Write your migrate up statements here

ALTER TABLE movie_stars ADD COLUMN hero_name VARCHAR(100);

---- create above / drop below ----

ALTER TABLE movie_stars DROP COLUMN hero_name;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
