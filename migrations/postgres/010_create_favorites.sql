-- Создание таблицы favorites
CREATE TABLE IF NOT EXISTS favorites (
    id SERIAL PRIMARY KEY,
    profile_id INTEGER NOT NULL,
    object_type VARCHAR(50) NOT NULL,
    object_id INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_favorites_profile FOREIGN KEY (profile_id) REFERENCES profiles(user_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_favorites_profile_id ON favorites(profile_id);

