package service

import "music-lib/internal/repository"

type LyricsService struct {
	lyricsRepository repository.ILyricsRepository
}

func NewLyricsService(lyrics repository.ILyricsRepository) *LyricsService {
	return &LyricsService{
		lyricsRepository: lyrics,
	}
}