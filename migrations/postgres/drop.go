package postgres

import (
	"database/sql"
	"fmt"
)

func DropTables(db *sql.DB) error {
	tables := []string{
		"collection_items",
		"collections",
		"favorites",
		"histories",
		"profiles",
		"resource_permissions",
		"song_genres",
		"couplets",
		"lyrics",
		"songs",
		"albums",
		"artists",
		"genres",
		"users",
	}

	for _, table := range tables {
		query := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table)
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}
	return nil
}
