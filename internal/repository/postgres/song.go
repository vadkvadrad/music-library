package postgres

import (
	"music-lib/pkg/db"

)

type SongRepository struct {
	db *db.Db
}

func NewSongRepository(db *db.Db) *SongRepository {
	return &SongRepository{
		db: db,
	}
}

