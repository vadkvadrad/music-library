package service

import (
	"errors"
	"music-lib/internal/dto/request"
	"music-lib/internal/model"
	"music-lib/internal/repository"
	"music-lib/pkg/er"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type ArtistService struct {
	artistRepository repository.IArtistRepository

	logger *zap.SugaredLogger
}

func NewArtistService(artist repository.IArtistRepository, log *zap.SugaredLogger) *ArtistService {
	return &ArtistService{
		artistRepository: artist,
		logger: log,
	}
}

func (s *ArtistService) NewArtist(ctx *gin.Context, body request.NewArtistRequest, userID uint) (*model.Artist, error) {
	_, err := s.artistRepository.GetByUserID(ctx, userID)
	if err == nil {
		return nil, er.ErrArtistLinked
	}

	exists := s.artistRepository.IsExists(ctx, body.ArtistName)
	if exists {
		return nil, er.ErrArtistExists
	}

	formationDate, err := time.Parse("2006-01-02", body.FormationYear)
	if err != nil {
		return nil, er.ErrDateFormat
	}

	return s.artistRepository.Create(ctx, &model.Artist{
		Name:          body.ArtistName,
		Description:   body.Description,
		FormationYear: formationDate,
		UserID:        userID,
	})
}


func (s *ArtistService) GetArtist(ctx *gin.Context, strID string) (*model.Artist, error) {
	id, err := strconv.Atoi(strID)
	if err != nil {
		return nil, &er.ValidationError{Message: err.Error()}
	}

	artist, err := s.artistRepository.GetWithAlbums(ctx, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, er.ErrArtistNotExists
		}
		return nil, &er.InternalError{Message: err.Error()}
	}

	return artist, nil
}


func (s *ArtistService) UpdateArtist(ctx *gin.Context, id uint, req request.UpdateArtistRequest) (*model.Artist, error) {
	s.logger.Debugw("Attempting to get artist",
		"id", id,
	)
	artist, err := s.artistRepository.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, er.ErrArtistNotExists
		}
		return nil, &er.InternalError{Message: err.Error()}
	}

	var formationDate time.Time
	if req.FormationYear != ""{
		formationDate, err = time.Parse("2006-01-02", req.FormationYear)
		if err != nil {
			return nil, er.ErrDateFormat
		}
	}
	
	s.logger.Debugw("Changing artist params",
		"Previous name", artist.Name,
		"Updated name", req.ArtistName,
		"Previous description", artist.Description,
		"Updated description", req.Description,
		"Previous formation year", artist.FormationYear,
		"Updated formation year", req.FormationYear,
	)
	artist.Name = req.ArtistName
	artist.Description = req.Description
	artist.FormationYear = formationDate

	s.logger.Debug("Artist updated successfully")
	return s.artistRepository.Update(ctx, artist)
}