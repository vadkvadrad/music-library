package postgres

import (
	"music-lib/pkg/db"
)

type LyricsRepository struct {
	db *db.Db
}

func NewLyricsRepository(db *db.Db) *LyricsRepository {
	return &LyricsRepository{
		db: db,
	}
}

