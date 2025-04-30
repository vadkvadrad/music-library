package service

import (
	"errors"
	"music-lib/internal/model"
	"music-lib/internal/repository"
	"music-lib/pkg/er"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
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


func (s *GenreService) UpdateGenre(ctx *gin.Context, id uint, nameToUpdate string) error {
	genre, err := s.genreRepo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Debugw("Can't find genre",
				"genre id", id,
				"error", err.Error(),
			)
			return er.ErrGenreNotExists
		}
		s.logger.Errorw("Error while finding genre",
			"error type", "internal",
			"error", err.Error(),
		)
		return &er.InternalError{Message: err.Error()}
	}

	genre.Name = nameToUpdate
	err = s.genreRepo.Update(ctx, genre)
	if err != nil {
		s.logger.Errorw("Error while updating genre",
			"error type", "internal",
			"error", err.Error(),
		)
		return &er.InternalError{Message: err.Error()}
	}

	s.logger.Debug("genre updated successfully")
	return nil
}