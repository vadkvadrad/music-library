package service

import (
	"errors"
	"fmt"
	"music-lib/internal/dto/request"
	"music-lib/internal/model"
	"music-lib/internal/repository"
	"music-lib/pkg/er"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AlbumService struct {
	albumRepository  repository.IAlbumRepository
	artistRepository repository.IArtistRepository
}

func NewAlbumService(album repository.IAlbumRepository, artist repository.IArtistRepository) *AlbumService {
	return &AlbumService{
		artistRepository: artist,
		albumRepository:  album,
	}
}

func (s *AlbumService) NewAlbum(ctx *gin.Context, body request.NewAlbumRequest, userID uint) error {
	artist, err := s.artistRepository.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return er.ErrArtistNotExists
		}
		return &er.InternalError{Message: fmt.Sprintf("NewAlbum: can't get artist: %s", err.Error())}
	}

	formationDate, err := time.Parse("2006-01-02", body.ReleaseDate)
	if err != nil {
		return er.ErrDateFormat
	}

	for _, album := range artist.Albums {
		if album.Title == body.Title {
			return er.ErrAlbumExists
		}
	}

	err = s.albumRepository.Create(ctx, &model.Album{
		Title:       body.Title,
		ArtistID:    artist.ID,
		Songs:       nil,
		ReleaseDate: formationDate,
		CoverArtURL: body.CoverArtURL,
	})

	if err != nil {
		return &er.InternalError{Message: fmt.Sprintf("NewAlbum: can't create album: %s", err.Error())}
	}
	return nil
}

func (s *AlbumService) GetAlbum(ctx *gin.Context, strID string) (*model.Album, error) {
	id, err := strconv.Atoi(strID)
	if err != nil {
		return nil, &er.ValidationError{Message: err.Error()}
	}

	album, err := s.albumRepository.GetWithSongs(ctx, uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, er.ErrAlbumNotExists
		}

		return nil, &er.InternalError{Message: err.Error()}
	}

	return album, nil
}
