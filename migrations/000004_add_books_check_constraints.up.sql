ALTER TABLE books ADD CONSTRAINT books_pages_check CHECK (pages >= 1);
ALTER TABLE books ADD CONSTRAINT genres_length_check CHECK (array_length(genres, 1) BETWEEN 1 AND 5);