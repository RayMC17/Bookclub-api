CREATE TABLE IF NOT EXISTS comments (
    id bigserial PRIMARY KEY, 
    created_at timestamp(0) WITH TIME ZONE NOT FULL DEFAULT NOW(),
    content text NOT NULL,
    author text NOT NULL,
    version integer NOT NULL DEFUALT 1
);

ALTER TABLE comments
DROP COLUMN IF EXISTS likes; 