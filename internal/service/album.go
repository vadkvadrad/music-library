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

func (s *AlbumService) NewAlbum(ctx *gin.Context, body request.NewAlbumRequest, userID uint) (*model.Album, error) {
	artist, err := s.artistRepository.GetByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, er.ErrArtistNotExists
		}
		return nil, &er.InternalError{Message: fmt.Sprintf("NewAlbum: can't get artist: %s", err.Error())}
	}

	formationDate, err := time.Parse("2006-01-02", body.ReleaseDate)
	if err != nil {
		return nil, er.ErrDateFormat
	}

	for _, album := range artist.Albums {
		if album.Title == body.Title {
			return nil, er.ErrAlbumExists
		}
	}

	album, err := s.albumRepository.Create(ctx, &model.Album{
		Title:       body.Title,
		ArtistID:    artist.ID,
		Songs:       nil,
		ReleaseDate: formationDate,
		CoverArtURL: body.CoverArtURL,
	})

	if err != nil {
		return nil, &er.InternalError{Message: fmt.Sprintf("NewAlbum: can't create album: %s", err.Error())}
	}
	return album, nil
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


func (s *AlbumService) GetArtistAlbum(ctx *gin.Context, userID uint, albumID uint) (*model.Album, error) {
	album, count, err := s.artistRepository.GetArtistAlbumByUserID(ctx, userID, albumID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, er.ErrAlbumNotExists
		}
		return nil, &er.InternalError{Message: err.Error()}
	}

	if count <= 0 {
		return nil, er.ErrAlbumNotExists
	}
	return album, nil
}
