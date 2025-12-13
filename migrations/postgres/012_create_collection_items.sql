-- Создание таблицы collection_items
CREATE TABLE IF NOT EXISTS collection_items (
    collection_id INTEGER NOT NULL,
    song_id INTEGER NOT NULL,
    position INTEGER NOT NULL,
    PRIMARY KEY (collection_id, song_id),
    CONSTRAINT fk_collection_items_collection FOREIGN KEY (collection_id) REFERENCES collections(id) ON DELETE CASCADE,
    CONSTRAINT fk_collection_items_song FOREIGN KEY (song_id) REFERENCES songs(id) ON DELETE CASCADE
);

