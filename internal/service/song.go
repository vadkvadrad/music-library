package service

import (
	"context"
	"music-lib/internal/dto/request"
	"music-lib/internal/model"
	"music-lib/internal/repository"
)

type SongService struct {
	songRepo repository.ISongRepository
	albumRepo repository.IAlbumRepository
}

func NewSongService(song repository.ISongRepository, album repository.IAlbumRepository) *SongService {
	return &SongService{
		songRepo: song,
		albumRepo: album,
	}
}

func (s *SongService) AddSong(ctx *context.Context, album *model.Album, song *request.NewSongRequest) error {
	
}