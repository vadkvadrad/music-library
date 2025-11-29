-- Создание таблицы couplets
CREATE TABLE IF NOT EXISTS couplets (
    id SERIAL PRIMARY KEY,
    lyrics_id INTEGER NOT NULL,
    number INTEGER NOT NULL,
    text TEXT NOT NULL,
    CONSTRAINT fk_couplets_lyrics FOREIGN KEY (lyrics_id) REFERENCES lyrics(song_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_couplets_lyrics_id ON couplets(lyrics_id);

