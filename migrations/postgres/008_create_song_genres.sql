-- Создание таблицы song_genres
CREATE TABLE IF NOT EXISTS song_genres (
    id SERIAL PRIMARY KEY,
    song_id INTEGER NOT NULL,
    genre_id INTEGER NOT NULL,
    CONSTRAINT fk_song_genres_song FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE,
    CONSTRAINT fk_song_genres_genre FOREIGN KEY (genre_id) REFERENCES genres(id) ON DELETE CASCADE,
    CONSTRAINT unique_song_genre UNIQUE (song_id, genre_id)
);

CREATE INDEX IF NOT EXISTS idx_song_genres_song_id ON song_genres(song_id);
CREATE INDEX IF NOT EXISTS idx_song_genres_genre_id ON song_genres(genre_id);

