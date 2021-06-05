-- Write your migrate up statements here
CREATE TABLE post (
    id SERIAL PRIMARY KEY,
    link TEXT,
    title TEXT,
    titletsv tsvector GENERATED ALWAYS AS (to_tsvector('russian', title)) STORED,
    description TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_post_link_unique ON post (link);
CREATE INDEX idx_post_titletsv_gin ON post USING GIN(titletsv);

---- create above / drop below ----

DROP TABLE post;

DROP INDEX IF EXISTS idx_post_link_unique;
DROP INDEX IF EXISTS idx_post_titletsv_gin; 

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
