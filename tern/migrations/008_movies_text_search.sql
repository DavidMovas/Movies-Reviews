ALTER TABLE movies ADD COLUMN search_vector TSVECTOR;

CREATE FUNCTION movies_search_vector_trigger() RETURNS TRIGGER AS $$
BEGIN
    new.search_vector :=
        setweight(to_tsvector('english', new.title), 'A') ||
        setweight(to_tsvector('english', new.description), 'B');
    return new;
END
$$ LANGUAGE plpgsql;

CREATE TRIGGER movies_search_vector_trigger
    BEFORE INSERT OR UPDATE ON movies
    FOR EACH ROW
EXECUTE PROCEDURE movies_search_vector_trigger();

CREATE INDEX idx_movies_search_vector ON movies USING GIN (search_vector);
---- create above / drop below ----
DROP INDEX idx_movies_search_vector;
DROP TRIGGER movies_search_vector_trigger ON movies;
DROP FUNCTION movies_search_vector_trigger();
ALTER TABLE movies DROP COLUMN search_vector;
