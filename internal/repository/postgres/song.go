package postgres

import (
	"context"
	"fmt"
	"music-lib/internal/model"
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

func (r *SongRepository) Create(ctx context.Context, entity *model.Song) (*model.Song, error) {
	err := r.db.WithContext(ctx).Create(entity).Error
	return entity, err
}

func (r *SongRepository) Update(ctx context.Context, entity *model.Song) (*model.Song, error) {
	err := r.db.WithContext(ctx).Save(entity).Error
	return entity, err
}

func (r *SongRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Song{}, id).Error
}

func (r *SongRepository) Search(ctx context.Context, query string, limit, offset int) ([]model.Song, int64, error) {
	var songs []model.Song
	db := r.db.WithContext(ctx).Model(&model.Song{})

	if query != "" {
		db = db.Where("name LIKE ?", "%"+query+"%")
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("error counting songs: %w", err)
	}

	err := db.Limit(limit).
		Offset(offset).
		Order("name ASC").
		Find(&songs).Error

	return songs, total, err
}

func (r *SongRepository) ExistsInAlbum(ctx context.Context, albumID uint, songName string) bool {
	var count int64
    err := r.db.WithContext(ctx).
        Model(&model.Song{}).
		Where("album_id = ?", albumID).
        Where("LOWER(title) = LOWER(?)", songName). 
        Limit(1).
        Count(&count).
        Error

    if err != nil {
        // Логирование ошибки 
        return false
    }

    return count > 0
}


func (r *SongRepository) GetByID(ctx context.Context, id uint) (*model.Song, error) {
	panic("SongRepository Implement GetByID")
}
func (r *SongRepository) GetByArtistID(ctx context.Context, artistID uint, sort string, limit, offset int) ([]model.Song, int64, error){
	panic("SongRepository Implement GetByArtistID")
}
func (r *SongRepository) GetByAlbumID(ctx context.Context, albumID uint, sort string, limit, offset int) ([]model.Song, int64, error){
	panic("SongRepository Implement GetByAlbumID")
}
func (r *SongRepository) GetFullInfo(ctx context.Context, id uint) (*model.Song, *model.Artist, *model.Album, error) {
	panic("SongRepository Implement GetFullInfo")
}