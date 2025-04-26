package postgres

import (
	"music-lib/internal/model"

	"gorm.io/gorm"
)

func MigrateTables(db *gorm.DB) error {
	return db.AutoMigrate(
		// User
		&model.User{},
		// Music
		&model.Artist{},
		&model.Album{},
		&model.Song{},
		&model.Lyrics{},
		&model.Couplet{},
		&model.Genre{},
		&model.SongGenre{},
		// Profile
		&model.Profile{},
		&model.Favorite{},
		&model.Collection{},
		&model.CollectionItem{},
		&model.History{},
	)
}
