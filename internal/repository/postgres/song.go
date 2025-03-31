package postgres

import (
	"music-lib/internal/model"
	"music-lib/pkg/db"

	"gorm.io/gorm/clause"
)

type SongRepository struct {
	Db *db.Db
}

func NewSongRepository(db *db.Db) *SongRepository {
	return &SongRepository{
		Db: db,
	}
}

func (repo *SongRepository) Create(song *model.Song) (*model.Song, error) {
	result := repo.Db.Create(song)
	if result.Error != nil {
		return nil, result.Error
	}
	return song, nil
}

func (repo *SongRepository) Update(song *model.Song) (*model.Song, error) {
	result := repo.Db.Clauses(clause.Returning{}).Updates(song)
	if result.Error != nil {
		return nil, result.Error
	}
	return song, nil
}

func (repo *SongRepository) Delete(id uint) error {
	result := repo.Db.Delete(&model.Song{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *SongRepository) FindByID(id uint) (*model.Song, error) {
	var song model.Song
	result := repo.Db.First(&song, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &song, nil
}

func (repo *SongRepository) Find(song, group string) (*model.Song, error) {
	var songRecord model.Song
	err := repo.Db.Where(`song = ? AND "group" = ?`, song, group).First(&songRecord).Error
	if err != nil {
		return nil, err
	}
	return &songRecord, nil
}


func (repo *SongRepository) FindByGroup(group string, limit, offset int) []model.Song {
	var songs []model.Song

	repo.Db.
		Table("songs").
		Where(`deleted_at is null AND "group" = ?`, group).
		Order("id asc").
		Limit(limit).
		Offset(offset). 
		Scan(&songs)
	return songs
}
