package service

import (
	"context"
	"errors"
	"music-lib/internal/dto/request"
	"music-lib/internal/model"
	"music-lib/internal/repository"
	"music-lib/pkg/er"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SongService struct {
	songRepo       repository.ISongRepository
	albumRepo      repository.IAlbumRepository
	songGenreRepo  repository.ISongGenreRepository
	genreRepo      repository.IGenreRepository
	lyricsRepo     repository.ILyricsRepository

	logger *zap.SugaredLogger
}

func NewSongService(
	song repository.ISongRepository,
	album repository.IAlbumRepository,
	songGenre repository.ISongGenreRepository,
	genre repository.IGenreRepository,
	lyrics repository.ILyricsRepository,
	sugar *zap.SugaredLogger,
) *SongService {
	return &SongService{
		songRepo:      song,
		albumRepo:     album,
		songGenreRepo: songGenre,
		genreRepo:     genre,
		lyricsRepo:    lyrics,
		logger:        sugar,
	}
}

func (s *SongService) AddSong(ctx context.Context, album *model.Album, songReq request.NewSongRequest) (*model.Song, error) {
	s.logger.Debugw("Attempting to add song",
		"album_id", album.ID,
		"artist_id", album.ArtistID,
		"song_title", songReq.Title,
		"duration", songReq.Duration,
	)
	if s.songRepo.ExistsInAlbum(ctx, album.ID, songReq.Title) {
		s.logger.Debugw("Song already exists in album",
			"album_id", album.ID,
			"song_title", songReq.Title,
			"error", er.ErrSongExists.Error(),
		)
		return nil,  er.ErrSongExists
	}

	s.logger.Debug("Attempting to create song")

	song, err := s.songRepo.Create(ctx, &model.Song{
		Title:      songReq.Title,
		AlbumID:    album.ID,
		ArtistID:   album.ArtistID,
		SongGenres: nil,
		Duration:   songReq.Duration,
		FilePath:   songReq.FilePath,
	})

	if err != nil {
		s.logger.Errorw("Failed to create song",
			"album_id", album.ID,
			"song_title", songReq.Title,
			"error", err.Error(),
		)
		return nil, &er.InternalError{Message: err.Error()}
	}

	s.logger.Debug("Attempting to add genres")
	err = s.addGenres(ctx, song.ID, songReq.Genres)
	if err != nil {
		s.logger.Errorw("Failed to add genres",
			"album_id", album.ID,
			"song_title", songReq.Title,
			"error", err.Error(),
		)
		return nil, &er.InternalError{Message: err.Error()}
	}

	s.logger.Debug("Attempting to add lyrics")
	err = s.addLyrics(ctx, song.ID, songReq.Lyrics)
	if err != nil {
		s.logger.Errorw("Failed to add lyrics",
			"album_id", album.ID,
			"song_title", songReq.Title,
			"error", err.Error(),
		)
		return nil, &er.InternalError{Message: err.Error()}
	}

	s.logger.Debug("Song created successfully")

	return song, nil
}

func (s *SongService) addGenres(ctx context.Context, songID uint, req []request.Genres) error {
	var genreIds []uint
	for _, genre := range req {
		genreIds = append(genreIds, genre.GenreID)
	}
	genres, err := s.genreRepo.GetByIds(ctx, genreIds)
	if err != nil {
		return &er.InternalError{Message: err.Error()}
	}

	for _, genre := range genres {
		_, err := s.songGenreRepo.Create(ctx, &model.SongGenre{
			SongID:  songID,
			GenreID: genre.ID,
		})
		if err != nil {
			return &er.InternalError{Message: err.Error()}
		}
	}
	return nil
}

func (s *SongService) addLyrics(ctx context.Context, songID uint, req request.AddLyrics) error {
	var lyrics model.Lyrics
	lyrics.SongID = songID

	for number, couplet := range req.Text {
		lyrics.Couplets = append(lyrics.Couplets, model.Couplet{
			LyricsID: songID,
			Number:   uint(number),
			Text:     couplet.Text,
		})
	}

	return s.lyricsRepo.Upsert(ctx, &lyrics)
}


func (s *SongService) GetSong(ctx context.Context, songID uint) (*model.Song, error) {
	song, err := s.songRepo.GetByID(ctx, songID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, er.ErrSongNotExists
		}
		return nil, &er.InternalError{Message: err.Error()}
	}

	return song, nil
}
