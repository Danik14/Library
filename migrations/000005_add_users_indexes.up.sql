CREATE INDEX IF NOT EXISTS users_firstName_idx ON users USING GIN (to_tsvector('simple', firstName));
CREATE INDEX IF NOT EXISTS users_lastName_idx ON users USING GIN (to_tsvector('simple', lastName));
CREATE INDEX IF NOT EXISTS users_email_idx ON users USING GIN (to_tsvector('simple', email));