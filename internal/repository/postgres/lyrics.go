package postgres

import (
	"context"
	"music-lib/internal/model"
	"music-lib/pkg/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type LyricsRepository struct {
	db *db.Db
}

func NewLyricsRepository(db *db.Db) *LyricsRepository {
	return &LyricsRepository{
		db: db,
	}
}

func (r *LyricsRepository) GetBySongID(ctx context.Context, songID uint) (*model.Lyrics, error) {
	var lyrics model.Lyrics
	err := r.db.WithContext(ctx).
		Preload("Couplets").
		Where("song_id = ?", songID).
		First(&lyrics).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &lyrics, nil
}

func (r *LyricsRepository) Upsert(ctx context.Context, lyrics *model.Lyrics) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "song_id"}},
			UpdateAll: true,
		}).
		Create(&lyrics).Error
}

func (r *LyricsRepository) DeleteBySongID(ctx context.Context, songID uint) error {
	result := r.db.WithContext(ctx).
		Where("song_id = ?", songID).
		Delete(&model.Lyrics{})
	
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}