package service

import (
	"fmt"
	"music-lib/internal/dto/request"
	"music-lib/internal/model"
	"music-lib/internal/repository"
	"music-lib/pkg/er"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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
		if err.Error() == "artist not found" {
			return nil, er.ErrArtistNotExists
		}
		return nil, &er.InternalError{Message: fmt.Sprintf("NewAlbum: can't get artist: %s", err.Error())}
	}

	formationDate, err := time.Parse("2006-01-02", body.ReleaseDate)
	if err != nil {
		return nil, er.ErrDateFormat
	}

	album, err := s.albumRepository.Create(ctx, &model.Album{
		Title:       body.Title,
		ArtistID:    artist.ID,
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
		if err.Error() == "album not found" {
			return nil, er.ErrAlbumNotExists
		}

		return nil, &er.InternalError{Message: err.Error()}
	}

	return album, nil
}

func (s *AlbumService) GetArtistAlbum(ctx *gin.Context, userID uint, albumID uint) (*model.Album, error) {
	// Сначала проверяем, есть ли у пользователя артист
	artist, err := s.artistRepository.GetByUserID(ctx, userID)
	if err != nil {
		if err.Error() == "artist not found" {
			return nil, er.ErrArtistNotExists
		}
		return nil, &er.InternalError{Message: fmt.Sprintf("GetArtistAlbum: can't get artist: %s", err.Error())}
	}

	// Получаем альбом
	album, err := s.albumRepository.GetByID(ctx, albumID)
	if err != nil {
		if err.Error() == "album not found" {
			return nil, er.ErrAlbumNotExists
		}
		return nil, &er.InternalError{Message: err.Error()}
	}

	// Проверяем, что альбом принадлежит артисту пользователя
	if album.ArtistID != artist.ID {
		return nil, &er.InternalError{Message: fmt.Sprintf("GetArtistAlbum: album %d does not belong to artist %d (user's artist is %d)", albumID, album.ArtistID, artist.ID)}
	}

	return album, nil
}

func (s *AlbumService) GetSongs(ctx *gin.Context, albumID uint) ([]model.Song, error) {
	songs, err := s.albumRepository.GetSongsByAlbumID(ctx, albumID)
	if err != nil {
		return nil, &er.InternalError{Message: err.Error()}
	}
	return songs, nil
}
