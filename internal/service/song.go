package service

import (
	"music-lib/internal/model"
	"music-lib/internal/repository"
	"music-lib/pkg/er"
)

type SongService struct {
	SongRepository repository.ISongRepository
}

func NewSongService(songRepository repository.ISongRepository) *SongService {
	return &SongService{
		SongRepository: songRepository,
	}
}

func (service *SongService) Add(song *model.Song) (*model.Song, error) {
	return service.SongRepository.Create(song)
}

func (service *SongService) Update(song *model.Song) (*model.Song, error) {
	return service.SongRepository.Update(song)
}

func (service *SongService) Delete(id uint) error {
	return service.SongRepository.Delete(id)
}

func (service *SongService) GetAndCompare(id uint, email string) (*model.Song, error) {
	song, err := service.SongRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	if song.Owner != email {
		return nil, er.ErrWrongUserCredentials
	}

	return song, nil
}

func (service *SongService) Get(song, group string) (*model.Song, error) {
	return service.SongRepository.Find(song, group)
}

func (service *SongService) GetByGroup(group string, limit, offset int) []model.Song {
	return service.SongRepository.FindByGroup(group, limit, offset)
}
