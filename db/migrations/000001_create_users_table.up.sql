CREATE TABLE IF NOT EXISTS users(
    id  UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    username        VARCHAR(25) UNIQUE,
    email           VARCHAR(40) UNIQUE,
    passwordHash    CHAR(60)    NOT NULL,
    created_at      TIMESTAMP DEFAULT now()
)
