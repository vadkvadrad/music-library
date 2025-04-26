package postgres

import (
	"context"
	"fmt"
	"music-lib/internal/model"
	"music-lib/pkg/db"

	"gorm.io/gorm"
)

type AlbumRepository struct {
	db *db.Db
}

func NewAlbumRepository(db *db.Db) *AlbumRepository {
	return &AlbumRepository{
		db: db,
	}
}

func (r *AlbumRepository) Create(ctx context.Context, entity *model.Album) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *AlbumRepository) Update(ctx context.Context, entity *model.Album) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *AlbumRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Album{}, id).Error
}

func (r *AlbumRepository) Search(ctx context.Context, query string, limit, offset int) ([]model.Album, int64, error) {
	var albums []model.Album
	db := r.db.WithContext(ctx).Model(&model.Album{})

	if query != "" {
		db = db.Where("name LIKE ?", "%"+query+"%")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting artists: %w", err)
	}

	err := db.Limit(limit).
		Preload("Songs").
		Offset(offset).
		Order("name ASC").
		Find(&albums).Error

	return albums, total, err
}


func (r *AlbumRepository) GetByID(ctx context.Context, id uint) (*model.Album, error) {
	var album *model.Album
	err := r.db.WithContext(ctx).
		First(&album, id).Error

	if err != nil {
		return nil, err
	}
	return album, nil
}


func (r *AlbumRepository) GetWithSongs(ctx context.Context, id uint) (*model.Album, error) {
	var album *model.Album
	err := r.db.WithContext(ctx).
		Preload("Songs", func(db *db.Db) *gorm.DB {
			return db.DB.Order("created_at DESC")
		}(r.db)).
		First(&album, id).Error

	if err != nil {
		return nil, err
	}

	return album, nil
}