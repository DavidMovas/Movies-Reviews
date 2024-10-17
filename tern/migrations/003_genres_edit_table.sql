-- Write your migrate up statements here

ALTER TABLE genres ADD CONSTRAINT unique_genre_name UNIQUE (name);

---- create above / drop below ----

ALTER TABLE genres DROP CONSTRAINT unique_genre_name;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
