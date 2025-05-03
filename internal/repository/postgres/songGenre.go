package postgres

import (
	"context"
	"music-lib/internal/model"
	"music-lib/pkg/db"
)

type SongGenreRepository struct {
	db *db.Db
}

func NewSongGenreRepository(db *db.Db) *SongGenreRepository {
	return &SongGenreRepository{
		db: db,
	}
}

func (r *SongGenreRepository) Create(ctx context.Context, entity *model.SongGenre) (*model.SongGenre, error) {
	err := r.db.WithContext(ctx).Create(entity).Error
	return entity, err
}

func (r *SongGenreRepository) Update(ctx context.Context, entity *model.SongGenre) (*model.SongGenre, error) {
	err := r.db.WithContext(ctx).Save(entity).Error
	return entity, err
}

func (r *SongGenreRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.SongGenre{}, id).Error
}
