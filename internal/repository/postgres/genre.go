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

func (r *GenreRepository) Create(ctx context.Context, entity *model.Genre) (*model.Genre, error) {
	err := r.db.WithContext(ctx).Create(entity).Error
	return entity, err
}

func (r *GenreRepository) Update(ctx context.Context, entity *model.Genre) (*model.Genre, error) {
	err := r.db.WithContext(ctx).Save(entity).Error
	return entity, err
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

func (r *GenreRepository) GetById(ctx context.Context, id uint) (*model.Genre, error) {
	var genre *model.Genre
	err := r.db.WithContext(ctx).
		First(&genre, id).Error

	if err != nil {
		return nil, err
	}
	return genre, nil
}

func (r *GenreRepository) GetByIds(ctx context.Context, ids []uint) ([]model.Genre, error) {
    if len(ids) == 0 {
        return []model.Genre{}, nil
    }

    var genres []model.Genre
    
    err := r.db.WithContext(ctx).
        Where("id IN ?", ids). // Условие IN для выборки по списку
        Find(&genres).Error
    
    if err != nil {
        return nil, err
    }
    
    return genres, nil
}
