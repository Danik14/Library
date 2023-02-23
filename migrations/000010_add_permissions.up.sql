CREATE TABLE IF NOT EXISTS permissions (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
code text NOT NULL
);
CREATE TABLE IF NOT EXISTS users_permissions (
user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE,
permission_id UUID NOT NULL REFERENCES permissions ON DELETE CASCADE,
PRIMARY KEY(user_id, permission_id)
);
-- Add the two permissions to the table.
INSERT INTO permissions (code)
VALUES
('books:read'),
('books:write'),
('users:read'),
('users:write'); 