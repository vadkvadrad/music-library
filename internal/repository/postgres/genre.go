package postgres

import (
	"context"
	"music-lib/internal/model"
	"music-lib/pkg/db"
)

type GenreRepository struct {
	db *db.Db
}

func NewGenreRepository(db *db.Db) *GenreRepository {
	return &GenreRepository{
		db: db,
	}
}

func (r *GenreRepository) Create(ctx context.Context, entity *model.Genre) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *GenreRepository) Update(ctx context.Context, entity *model.Genre) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *GenreRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Genre{}, id).Error
}

func (r *GenreRepository) IsExists(ctx context.Context, name string) bool {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.Genre{}).
		Where("LOWER(name) = LOWER(?)", name).
		Limit(1).
		Count(&count).
		Error

	if err != nil {
		// Логирование ошибки
		return false
	}

	return count > 0
}
