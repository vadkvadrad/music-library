-- Создание таблицы histories
CREATE TABLE IF NOT EXISTS histories (
    id SERIAL PRIMARY KEY,
    profile_id INTEGER NOT NULL,
    song_id INTEGER NOT NULL,
    played_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_histories_profile FOREIGN KEY (profile_id) REFERENCES profiles(user_id) ON DELETE CASCADE,
    CONSTRAINT fk_histories_song FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_histories_profile_id ON histories(profile_id);
CREATE INDEX IF NOT EXISTS idx_histories_song_id ON histories(song_id);
CREATE INDEX IF NOT EXISTS idx_histories_played_at ON histories(played_at);

