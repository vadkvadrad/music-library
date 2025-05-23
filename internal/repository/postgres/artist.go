package postgres

import (
	"context"
	"fmt"
	"music-lib/internal/model"
	"music-lib/pkg/db"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ArtistRepository struct {
	db *db.Db
}

func NewArtistRepository(db *db.Db) *ArtistRepository {
	return &ArtistRepository{
		db: db,
	}
}

func (r *ArtistRepository) Create(ctx context.Context, entity *model.Artist) (*model.Artist, error) {
	err := r.db.WithContext(ctx).Create(entity).Error
	return entity, err
}

func (r *ArtistRepository) Update(ctx context.Context, entity *model.Artist) (*model.Artist, error) {
	result := r.db.Clauses(clause.Returning{}).Updates(entity)
	if result.Error != nil {
		return nil, result.Error
	}
	return entity, nil
}

func (r *ArtistRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Artist{}, id).Error
}

func (r *ArtistRepository) Search(ctx context.Context, query string, limit, offset int) ([]model.Artist, int64, error) {
	var artists []model.Artist
	db := r.db.WithContext(ctx).Model(&model.Artist{})

	if query != "" {
		db = db.Where("LOWER(name) LIKE LOWER(?)", "%"+query+"%")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting artists: %w", err)
	}

	err := db.Limit(limit).
		Offset(offset).
		Order("name ASC").
		Find(&artists).Error

	return artists, total, err
}

// Реализация специфичных методов артиста
func (r *ArtistRepository) GetByID(ctx context.Context, id uint) (*model.Artist, error) {
	var artist *model.Artist
	err := r.db.WithContext(ctx).
		First(&artist, id).Error

	if err != nil {
		return nil, err
	}
	return artist, nil
}

func (r *ArtistRepository) GetByUserID(ctx context.Context, userID uint) (*model.Artist, error) {
	var artist *model.Artist
	err := r.db.WithContext(ctx).
		Preload("Albums", func(db *db.Db) *gorm.DB {
			return db.DB.Order("release_date DESC")
		}(r.db)).
		First(&artist, "user_id = ?", userID).Error

	if err != nil {
		return nil, err
	}
	return artist, nil
}

func (r *ArtistRepository) GetWithAlbums(ctx context.Context, id uint) (*model.Artist,  error) {
	var artist *model.Artist
	err := r.db.WithContext(ctx).
		Preload("Albums", func(db *db.Db) *gorm.DB {
			return db.DB.Order("release_date DESC")
		}(r.db)).
		First(&artist, id).Error

	if err != nil {
		return nil, err
	}

	return artist, nil
}

func (r *ArtistRepository) IsExists(ctx context.Context, name string) bool {
    var count int64
    err := r.db.WithContext(ctx).
        Model(&model.Artist{}).
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

func (r *ArtistRepository) GetArtistAlbumByUserID(ctx context.Context, userID uint, albumID uint) (*model.Album, int, error) {
	var artist *model.Artist
	err := r.db.WithContext(ctx).
		Preload("Albums", func(db *db.Db) *gorm.DB {
			return db.DB.Where("ID = ?", albumID)
		}(r.db)).
		First(&artist, "user_id = ?", userID).Error

	if err != nil {
		return nil, 0, err
	}
	if len(artist.Albums) <= 0 {
		return nil, 0, nil
	} 
	return &artist.Albums[0], len(artist.Albums), nil
}