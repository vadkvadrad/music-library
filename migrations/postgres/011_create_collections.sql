-- Создание таблицы collections
CREATE TABLE IF NOT EXISTS collections (
    id SERIAL PRIMARY KEY,
    profile_id INTEGER NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_collections_profile FOREIGN KEY (profile_id) REFERENCES profiles(user_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_collections_profile_id ON collections(profile_id);

