-- Write your migrate up statements here
ALTER TABLE users
ADD COLUMN avatar_url VARCHAR(255),
ADD COLUMN bio TEXT;

---- create above / drop below ----

ALTER TABLE users
DROP COLUMN avatar_url,
DROP COLUMN bio;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
