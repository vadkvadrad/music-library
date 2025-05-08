package service

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"music-lib/internal/dto/request"
	"music-lib/internal/model"
	"music-lib/internal/service/mocks"
	"music-lib/pkg/er"
)

func TestNewAlbum_ArtistNotExists(t *testing.T) {
	// Arrange
	mockArtistRepo := &mocks.MockArtistRepo{
		GetByUserIDFunc: func(ctx context.Context, userID uint) (*model.Artist, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	service := NewAlbumService(nil, mockArtistRepo)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := request.NewAlbumRequest{
		Title:       "Test Album",
		ReleaseDate: "2023-01-01",
	}

	// Act
	album, err := service.NewAlbum(ctx, req, 1)

	// Assert
	assert.Nil(t, album)
	assert.Equal(t, er.ErrArtistNotExists, err)
}

func TestNewAlbum_InvalidDateFormat(t *testing.T) {
	// Arrange
	mockArtistRepo := &mocks.MockArtistRepo{
		GetByUserIDFunc: func(ctx context.Context, userID uint) (*model.Artist, error) {
			return &model.Artist{}, nil
		},
	}
	service := NewAlbumService(nil, mockArtistRepo)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := request.NewAlbumRequest{
		Title:       "Test Album",
		ReleaseDate: "invalid-date",
	}

	// Act
	album, err := service.NewAlbum(ctx, req, 1)

	// Assert
	assert.Nil(t, album)
	assert.Equal(t, er.ErrDateFormat, err)
}

func TestNewAlbum_AlbumExists(t *testing.T) {
	// Arrange
	mockArtistRepo := &mocks.MockArtistRepo{
		GetByUserIDFunc: func(ctx context.Context, userID uint) (*model.Artist, error) {
			return &model.Artist{
				Albums: []model.Album{{Title: "Test Album"}},
			}, nil
		},
	}
	service := NewAlbumService(nil, mockArtistRepo)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := request.NewAlbumRequest{
		Title:       "Test Album",
		ReleaseDate: "2023-01-01",
	}

	// Act
	album, err := service.NewAlbum(ctx, req, 1)

	// Assert
	assert.Nil(t, album)
	assert.Equal(t, er.ErrAlbumExists, err)
}

func TestNewAlbum_InternalErrorOnArtistFetch(t *testing.T) {
	// Arrange
	mockArtistRepo := &mocks.MockArtistRepo{
		GetByUserIDFunc: func(ctx context.Context, userID uint) (*model.Artist, error) {
			return nil, errors.New("database error")
		},
	}
	service := NewAlbumService(nil, mockArtistRepo)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := request.NewAlbumRequest{
		Title:       "Test Album",
		ReleaseDate: "2023-01-01",
	}

	// Act
	album, err := service.NewAlbum(ctx, req, 1)

	// Assert
	assert.Nil(t, album)
	assert.IsType(t, &er.InternalError{}, err)
	assert.Equal(t, "NewAlbum: can't get artist: database error", err.(*er.InternalError).Message)
}

func TestNewAlbum_Success(t *testing.T) {
	// Arrange
	mockArtistRepo := &mocks.MockArtistRepo{
		GetByUserIDFunc: func(ctx context.Context, userID uint) (*model.Artist, error) {
			return &model.Artist{ID: 1}, nil
		},
	}
	mockAlbumRepo := &mocks.MockAlbumRepo{
		CreateFunc: func(ctx context.Context, entity *model.Album) (*model.Album, error) {
			return &model.Album{
				ID:          1,
				Title:       entity.Title,
				ArtistID:    entity.ArtistID,
				ReleaseDate: entity.ReleaseDate,
			}, nil
		},
	}
	service := NewAlbumService(mockAlbumRepo, mockArtistRepo)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := request.NewAlbumRequest{
		Title:       "Test Album",
		ReleaseDate: "2023-01-01",
	}
	expectedDate, _ := time.Parse("2006-01-02", "2023-01-01")

	// Act
	album, err := service.NewAlbum(ctx, req, 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, album)
	assert.Equal(t, uint(1), album.ID)
	assert.Equal(t, "Test Album", album.Title)
	assert.Equal(t, uint(1), album.ArtistID)
	assert.Equal(t, expectedDate, album.ReleaseDate)
}

func TestGetAlbum_InvalidID(t *testing.T) {
	// Arrange
	mockAlbumRepo := &mocks.MockAlbumRepo{}
	service := NewAlbumService(mockAlbumRepo, nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Act
	album, err := service.GetAlbum(ctx, "invalid-id")

	// Assert
	assert.Nil(t, album)
	assert.IsType(t, &er.ValidationError{}, err)
	assert.Contains(t, err.(*er.ValidationError).Message, "invalid syntax")
}

func TestGetAlbum_AlbumNotFound(t *testing.T) {
	// Arrange
	mockAlbumRepo := &mocks.MockAlbumRepo{
		GetWithSongsFunc: func(ctx context.Context, id uint) (*model.Album, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	service := NewAlbumService(mockAlbumRepo, nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Act
	album, err := service.GetAlbum(ctx, "1")

	// Assert
	assert.Nil(t, album)
	assert.Equal(t, er.ErrAlbumNotExists, err)
}

func TestGetAlbum_Success(t *testing.T) {
	// Arrange
	mockAlbumRepo := &mocks.MockAlbumRepo{
		GetWithSongsFunc: func(ctx context.Context, id uint) (*model.Album, error) {
			return &model.Album{
				ID:    1,
				Title: "Test Album",
			}, nil
		},
	}
	service := NewAlbumService(mockAlbumRepo, nil)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Act
	album, err := service.GetAlbum(ctx, "1")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, album)
	assert.Equal(t, uint(1), album.ID)
	assert.Equal(t, "Test Album", album.Title)
}

func TestGetArtistAlbum_AlbumNotFound(t *testing.T) {
	// Arrange
	mockArtistRepo := &mocks.MockArtistRepo{
		GetArtistAlbumByUserIDFunc: func(ctx context.Context, userID uint, albumID uint) (*model.Album, int, error) {
			return nil, 0, gorm.ErrRecordNotFound
		},
	}
	service := NewAlbumService(nil, mockArtistRepo)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Act
	album, err := service.GetArtistAlbum(ctx, 1, 1)

	// Assert
	assert.Nil(t, album)
	assert.Equal(t, er.ErrAlbumNotExists, err)
}

func TestGetArtistAlbum_NoAlbums(t *testing.T) {
	// Arrange
	mockArtistRepo := &mocks.MockArtistRepo{
		GetArtistAlbumByUserIDFunc: func(ctx context.Context, userID uint, albumID uint) (*model.Album, int, error) {
			return nil, 0, nil
		},
	}
	service := NewAlbumService(nil, mockArtistRepo)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Act
	album, err := service.GetArtistAlbum(ctx, 1, 1)

	// Assert
	assert.Nil(t, album)
	assert.Equal(t, er.ErrAlbumNotExists, err)
}

func TestGetArtistAlbum_Success(t *testing.T) {
	// Arrange
	mockArtistRepo := &mocks.MockArtistRepo{
		GetArtistAlbumByUserIDFunc: func(ctx context.Context, userID uint, albumID uint) (*model.Album, int, error) {
			return &model.Album{
				ID:    1,
				Title: "Test Album",
			}, 1, nil
		},
	}
	service := NewAlbumService(nil, mockArtistRepo)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Act
	album, err := service.GetArtistAlbum(ctx, 1, 1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, album)
	assert.Equal(t, uint(1), album.ID)
	assert.Equal(t, "Test Album", album.Title)
}