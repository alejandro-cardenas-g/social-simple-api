CREATE TABLE IF NOT EXISTS roles (
    id bigserial PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    level INT NOT NULL DEFAULT 0,
    description TEXT
);

INSERT INTO roles (name, description, level) 
VALUES (
    'user',
    'An user can create posts and comments',
    1
), (
    'moderator',
    'A moderator can update posts',
    2
),  (
    'admin',
    'An admin can update and delete posts',
    3
)

