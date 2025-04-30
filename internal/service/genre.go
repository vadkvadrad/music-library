package service

import (
	"music-lib/internal/model"
	"music-lib/internal/repository"
	"music-lib/pkg/er"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type GenreService struct {
	genreRepo repository.IGenreRepository

	logger *zap.SugaredLogger
}

func NewGenreService(genre repository.IGenreRepository, sugar *zap.SugaredLogger) *GenreService {
	return &GenreService{
		genreRepo: genre,
		logger: sugar,
	}
}


func (s *GenreService) NewGenre(ctx *gin.Context, genreName string) error {
	if s.genreRepo.IsExists(ctx, genreName) {
		s.logger.Debugw("Genre already exists",
			"genre name", genreName,
		)
		return er.ErrGenreExists
	}

	err := s.genreRepo.Create(ctx, &model.Genre{
		Name: genreName,
		SongGenres: nil,
	})

	if err != nil {
		s.logger.Errorw("Can't create genre",
			"genre name", genreName,
			"error type", "internal",
			"error", err.Error(),
		)
		return &er.InternalError{Message: err.Error()}
	}

	s.logger.Debugw("New genre added successfully",
		"genre name", genreName,
	)
	return nil
}