CREATE TABLE IF NOT EXISTS recipes(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recipe_name VARCHAR(20) NOT NULL,
    description TEXT    NOT NULL,
    owner_id    UUID  REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT now(),
    deleted_at      TIMESTAMP
);

CREATE TABLE IF NOT EXISTS ingredients(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(25)    NOT NULL,
    qty  DECIMAL NOT NULL CHECK ( qty>0 ),
    measurement VARCHAR(8) NOT NULL,
    recipe_id   uuid  REFERENCES recipes(id),
    CONSTRAINT unique_key
        UNIQUE(name,recipe_id)
)
