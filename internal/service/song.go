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
	songRepo      repository.ISongRepository
	albumRepo     repository.IAlbumRepository
	songGenreRepo repository.ISongGenreRepository
	genreRepo     repository.IGenreRepository
	lyricsRepo    repository.ILyricsRepository

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
		return nil, er.ErrSongExists
	}

	s.logger.Debug("Attempting to create song")

	albumID := album.ID
	song, err := s.songRepo.Create(ctx, &model.Song{
		Title:    songReq.Title,
		AlbumID:  &albumID,
		ArtistID: album.ArtistID,
		Duration: songReq.Duration,
		FilePath: songReq.FilePath,
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

func (s *SongService) GetSong(ctx context.Context, songID uint) (*model.Song, *model.Lyrics, []model.Genre, error) {
	song, err := s.songRepo.GetByID(ctx, songID)
	if err != nil {
		if err.Error() == "song not found" {
			return nil, nil, nil, er.ErrSongNotExists
		}
		return nil, nil, nil, &er.InternalError{Message: err.Error()}
	}

	// Загружаем лирику (может быть nil, если лирика не найдена)
	lyrics, err := s.lyricsRepo.GetBySongID(ctx, songID)
	if err != nil {
		// Если ошибка при загрузке лирики, продолжаем без лирики
		lyrics = nil
	}

	// Загружаем жанры песни
	var genres []model.Genre
	songGenres, err := s.songGenreRepo.GetBySongID(ctx, songID)
	if err == nil && len(songGenres) > 0 {
		genreIDs := make([]uint, len(songGenres))
		for i, sg := range songGenres {
			genreIDs[i] = sg.GenreID
		}
		genres, err = s.genreRepo.GetByIds(ctx, genreIDs)
		if err != nil {
			// Если ошибка при загрузке жанров, продолжаем без жанров
			genres = []model.Genre{}
		}
	}

	return song, lyrics, genres, nil
}

func (s *SongService) UpdateSong(ctx context.Context, songID uint, songReq request.UpdateSongRequest) (*model.Song, error) {
	s.logger.Debugw("Attempting to update song",
		"song_id", songID,
		"song_title", songReq.Title,
		"duration", songReq.Duration,
	)

	song, err := s.songRepo.GetByID(ctx, songID)
	if err != nil {
		if err.Error() == "song not found" {
			return nil, er.ErrSongNotExists
		}
		return nil, &er.InternalError{Message: err.Error()}
	}

	// Обновляем поля песни, если они предоставлены
	if songReq.Title != "" {
		song.Title = songReq.Title
	}
	if songReq.Duration != 0 {
		song.Duration = songReq.Duration
	}
	if songReq.FilePath != "" {
		song.FilePath = songReq.FilePath
	}

	updatedSong, err := s.songRepo.Update(ctx, song)
	if err != nil {
		s.logger.Errorw("Failed to update song",
			"song_id", songID,
			"error", err.Error(),
		)
		return nil, &er.InternalError{Message: err.Error()}
	}

	// Обновляем жанры, если они предоставлены
	if len(songReq.Genres) > 0 {
		// Удаляем старые жанры
		err = s.removeGenres(ctx, songID)
		if err != nil {
			s.logger.Errorw("Failed to remove old genres",
				"song_id", songID,
				"error", err.Error(),
			)
			return nil, &er.InternalError{Message: err.Error()}
		}

		// Добавляем новые жанры
		err = s.addGenres(ctx, songID, songReq.Genres)
		if err != nil {
			s.logger.Errorw("Failed to add genres",
				"song_id", songID,
				"error", err.Error(),
			)
			return nil, &er.InternalError{Message: err.Error()}
		}
	}

	// Обновляем лирику, если она предоставлена
	if len(songReq.Lyrics.Text) > 0 {
		err = s.addLyrics(ctx, songID, songReq.Lyrics)
		if err != nil {
			s.logger.Errorw("Failed to update lyrics",
				"song_id", songID,
				"error", err.Error(),
			)
			return nil, &er.InternalError{Message: err.Error()}
		}
	}

	s.logger.Debug("Song updated successfully")
	return updatedSong, nil
}

func (s *SongService) removeGenres(ctx context.Context, songID uint) error {
	// Удаляем все связи жанров с песней
	err := s.songGenreRepo.DeleteBySongID(ctx, songID)
	if err != nil {
		return &er.InternalError{Message: err.Error()}
	}
	return nil
}
