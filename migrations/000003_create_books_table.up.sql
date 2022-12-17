CREATE TABLE IF NOT EXISTS books (
id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,
created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
title text NOT NULL,
author text NOT NULL,
year integer NOT NULL,
pages integer NOT NULL,
genres text[] NOT NULL,
version integer NOT NULL DEFAULT 1
);