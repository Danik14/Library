ALTER TABLE books ADD COLUMN IF NOT EXISTS new_id UUID NULL;
UPDATE books SET new_id = CAST(LPAD(TO_HEX(id), 32, '0') AS UUID);
ALTER TABLE books DROP COLUMN IF EXISTS id;
ALTER TABLE books RENAME COLUMN new_id TO id;
ALTER TABLE books ALTER COLUMN id SET NOT NULL;