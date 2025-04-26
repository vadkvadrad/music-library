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
	"gorm.io/gorm"
)

type ArtistService struct {
	artistRepository repository.IArtistRepository
}

func NewArtistService(artist repository.IArtistRepository) *ArtistService {
	return &ArtistService{
		artistRepository: artist,
	}
}

func (s *ArtistService) NewArtist(ctx *gin.Context, body request.NewArtistRequest, userID uint) error {
	_, err := s.artistRepository.GetByUserID(ctx, userID)
	if err == nil {
		return er.ErrArtistLinked
	}

	exists := s.artistRepository.IsExists(ctx, body.ArtistName)
	if exists {
		return er.ErrArtistExists
	}

	formationDate, err := time.Parse("2006-01-02", body.FormationYear)
	if err != nil {
		return er.ErrDateFormat
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