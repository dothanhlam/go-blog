-- Add a tsvector column to store the document vector for the title.
ALTER TABLE posts ADD COLUMN title_tsv tsvector;

-- Create a function to automatically update the title_tsv column.
CREATE OR REPLACE FUNCTION update_posts_tsvector()
RETURNS TRIGGER AS $$
BEGIN
    NEW.title_tsv = to_tsvector('english', NEW.title);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create a trigger that calls the function on insert or update.
CREATE TRIGGER posts_tsvector_update
BEFORE INSERT OR UPDATE ON posts
FOR EACH ROW EXECUTE PROCEDURE update_posts_tsvector();

-- Create a GIN index for fast full-text search.
CREATE INDEX posts_title_tsv_idx ON posts USING GIN(title_tsv);