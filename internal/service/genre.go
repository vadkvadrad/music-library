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

func NewGenreService(
	genre repository.IGenreRepository,
	sugar *zap.SugaredLogger) *GenreService {
	return &GenreService{
		genreRepo: genre,
		logger:    sugar,
	}
}

func (s *GenreService) NewGenre(ctx *gin.Context, genreName string) error {
	if s.genreRepo.IsExists(ctx, genreName) {
		s.logger.Debugw("Genre already exists",
			"genre name", genreName,
		)
		return er.ErrGenreExists
	}

	_, err := s.genreRepo.Create(ctx, &model.Genre{
		Name: genreName,
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
		if err.Error() == "genre not found" {
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
	_, err = s.genreRepo.Update(ctx, genre)
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

func (s *GenreService) GetAllGenres(ctx *gin.Context) ([]model.Genre, error) {
	genres, err := s.genreRepo.GetAll(ctx)
	if err != nil {
		s.logger.Errorw("Error while getting all genres",
			"error type", "internal",
			"error", err.Error(),
		)
		return nil, &er.InternalError{Message: err.Error()}
	}
	return genres, nil
}

func (s *GenreService) GetGenre(ctx *gin.Context, id uint) (*model.Genre, error) {
	genre, err := s.genreRepo.GetById(ctx, id)
	if err != nil {
		if err.Error() == "genre not found" {
			return nil, er.ErrGenreNotExists
		}
		return nil, &er.InternalError{Message: err.Error()}
	}
	return genre, nil
}
