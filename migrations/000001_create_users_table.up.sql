CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    createdAt TIMESTAMP NOT NULL DEFAULT NOW(),
    firstName VARCHAR(30) NOT NULL,
    lastName VARCHAR(30) NOT NULL,
    email VARCHAR(100) NOT NULL,
    hashedPassword bytea NOT NULL,
    dob TIMESTAMP NOT NULL,
    version INTEGER NOT NULL DEFAULT 1
);