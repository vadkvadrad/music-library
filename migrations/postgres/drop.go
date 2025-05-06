package postgres

import (
    "gorm.io/gorm"
)

func DropTables(db *gorm.DB) error {
    tables := []string{
        "song_genres",
        "genres",
        "couplets",
        "lyrics",
        "songs",
        "albums",
        "artists",
        "histories",
        "collection_items",
        "collections",
        "favorites",
        "profiles",
        "users",
        "resource_permission",
    }

    for _, table := range tables {
        if err := db.Migrator().DropTable(table); err != nil {
            return err
        }
    }
    return nil
}