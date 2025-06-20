ALTER TABLE IF EXISTS users
ADD COLUMN role_id BIGINT REFERENCES roles(id) NULL;

UPDATE users 
SET role_id = (
    SELECT id FROM roles WHERE name = 'user'
);

ALTER TABLE users
ALTER COLUMN role_id SET NOT NULL;