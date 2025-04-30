package service

import (
	"context"
	"music-lib/internal/dto/request"
	"music-lib/internal/model"
	"music-lib/internal/repository"
	"music-lib/pkg/er"

	"go.uber.org/zap"
)

type SongService struct {
	songRepo repository.ISongRepository
	albumRepo repository.IAlbumRepository
	
	logger *zap.SugaredLogger
}

func NewSongService(song repository.ISongRepository, album repository.IAlbumRepository, sugar *zap.SugaredLogger) *SongService {
	return &SongService{
		songRepo: song,
		albumRepo: album,
		logger: sugar,
	}
}

func (s *SongService) AddSong(ctx context.Context, album *model.Album, song request.NewSongRequest) error {
	s.logger.Debugw("Attempting to add song",
        "album_id", album.ID,
        "artist_id", album.ArtistID,
        "song_title", song.Title,
        "duration", song.Duration,
    )
	if s.songRepo.ExistsInAlbum(ctx, album.ID, song.Title) {
	    s.logger.Debugw("Song already exists in album",
            "album_id", album.ID,
            "song_title", song.Title,
            "error", er.ErrSongExists.Error(),
        )
		return er.ErrSongExists
	}

	s.logger.Debug("Attempting to create song")

	err := s.songRepo.Create(ctx, &model.Song{
		Title: song.Title,
		AlbumID: album.ID,
		ArtistID: album.ArtistID,
		SongGenres: nil,
		Duration: song.Duration,
		FilePath: song.FilePath,
	})

	if err != nil {
		s.logger.Errorw("Failed to create song",
            "album_id", album.ID,
            "song_title", song.Title,
            "error", err.Error(),
        )
		return &er.InternalError{Message: err.Error()}
	}
	return nil
}