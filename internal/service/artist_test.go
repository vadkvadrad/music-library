package service

import (
	"context"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"music-lib/internal/dto/request"
	"music-lib/internal/model"
	"music-lib/internal/service/mocks"
	"music-lib/pkg/er"
)

func TestNewArtist_UserAlreadyLinked(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := &mocks.MockArtistRepo{
		GetByUserIDFunc: func(ctx context.Context, userID uint) (*model.Artist, error) {
			return &model.Artist{}, nil
		},
	}
	service := NewArtistService(mockRepo, logger)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := request.NewArtistRequest{
		ArtistName:    "Test Artist",
		Description:   "Description",
		FormationYear: "2020-01-01",
	}
	artist, err := service.NewArtist(ctx, req, 1)
	assert.Nil(t, artist)
	assert.Error(t, err)
	assert.Equal(t, er.ErrArtistLinked, err)
}

func TestNewArtist_ArtistExists(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := &mocks.MockArtistRepo{
		GetByUserIDFunc: func(ctx context.Context, userID uint) (*model.Artist, error) {
			return nil, gorm.ErrRecordNotFound
		},
		IsExistsFunc: func(ctx context.Context, name string) bool {
			return true
		},
	}
	service := NewArtistService(mockRepo, logger)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := request.NewArtistRequest{
		ArtistName:    "Test Artist",
		Description:   "Description",
		FormationYear: "2020-01-01",
	}
	artist, err := service.NewArtist(ctx, req, 1)
	assert.Nil(t, artist)
	assert.Error(t, err)
	assert.Equal(t, er.ErrArtistExists, err)
}

func TestNewArtist_InvalidDateFormat(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := &mocks.MockArtistRepo{
		GetByUserIDFunc: func(ctx context.Context, userID uint) (*model.Artist, error) {
			return nil, gorm.ErrRecordNotFound
		},
		IsExistsFunc: func(ctx context.Context, name string) bool {
			return false
		},
	}
	service := NewArtistService(mockRepo, logger)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := request.NewArtistRequest{
		ArtistName:    "Test Artist",
		Description:   "Description",
		FormationYear: "invalid-date",
	}
	artist, err := service.NewArtist(ctx, req, 1)
	assert.Nil(t, artist)
	assert.Error(t, err)
	assert.Equal(t, er.ErrDateFormat, err)
}

func TestNewArtist_Success(t *testing.T) {
	logger := zap.NewNop().Sugar()
	expectedArtist := &model.Artist{
		ID:            1,
		Name:          "Test Artist",
		Description:   "Description",
		FormationYear: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
		UserID:        1,
	}
	mockRepo := &mocks.MockArtistRepo{
		GetByUserIDFunc: func(ctx context.Context, userID uint) (*model.Artist, error) {
			return nil, gorm.ErrRecordNotFound
		},
		IsExistsFunc: func(ctx context.Context, name string) bool {
			return false
		},
		CreateFunc: func(ctx context.Context, entity *model.Artist) (*model.Artist, error) {
			return expectedArtist, nil
		},
	}
	service := NewArtistService(mockRepo, logger)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := request.NewArtistRequest{
		ArtistName:    "Test Artist",
		Description:   "Description",
		FormationYear: "2020-01-01",
	}
	artist, err := service.NewArtist(ctx, req, 1)
	assert.NoError(t, err)
	assert.NotNil(t, artist)
	assert.Equal(t, expectedArtist, artist)
}

func TestGetArtist_InvalidID(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := &mocks.MockArtistRepo{}
	service := NewArtistService(mockRepo, logger)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	artist, err := service.GetArtist(ctx, "abc")
	assert.Nil(t, artist)
	assert.Error(t, err)
	assert.IsType(t, &er.ValidationError{}, err)
}

func TestGetArtist_NotFound(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := &mocks.MockArtistRepo{
		GetWithAlbumsFunc: func(ctx context.Context, id uint) (*model.Artist, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	service := NewArtistService(mockRepo, logger)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	artist, err := service.GetArtist(ctx, "1")
	assert.Nil(t, artist)
	assert.Error(t, err)
	assert.Equal(t, er.ErrArtistNotExists, err)
}

func TestGetArtist_Success(t *testing.T) {
	logger := zap.NewNop().Sugar()
	expectedArtist := &model.Artist{
		ID: 1,
	}
	mockRepo := &mocks.MockArtistRepo{
		GetWithAlbumsFunc: func(ctx context.Context, id uint) (*model.Artist, error) {
			return expectedArtist, nil
		},
	}
	service := NewArtistService(mockRepo, logger)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	artist, err := service.GetArtist(ctx, "1")
	assert.NoError(t, err)
	assert.NotNil(t, artist)
	assert.Equal(t, expectedArtist, artist)
}

func TestUpdateArtist_NotFound(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := &mocks.MockArtistRepo{
		GetByIDFunc: func(ctx context.Context, id uint) (*model.Artist, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	service := NewArtistService(mockRepo, logger)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := request.UpdateArtistRequest{
		ArtistName:    "Updated Artist",
		Description:   "Updated Description",
		FormationYear: "2021-01-01",
	}
	artist, err := service.UpdateArtist(ctx, 1, req)
	assert.Nil(t, artist)
	assert.Error(t, err)
	assert.Equal(t, er.ErrArtistNotExists, err)
}

func TestUpdateArtist_InvalidDateFormat(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockRepo := &mocks.MockArtistRepo{
		GetByIDFunc: func(ctx context.Context, id uint) (*model.Artist, error) {
			return &model.Artist{}, nil
		},
	}
	service := NewArtistService(mockRepo, logger)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := request.UpdateArtistRequest{
		ArtistName:    "Updated Artist",
		Description:   "Updated Description",
		FormationYear: "invalid-date",
	}
	artist, err := service.UpdateArtist(ctx, 1, req)
	assert.Nil(t, artist)
	assert.Error(t, err)
	assert.Equal(t, er.ErrDateFormat, err)
}

func TestUpdateArtist_Success(t *testing.T) {
	logger := zap.NewNop().Sugar()
	originalArtist := &model.Artist{
		ID:            1,
		Name:          "Original Artist",
		Description:   "Original Description",
		FormationYear: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	updatedArtist := &model.Artist{
		ID:            1,
		Name:          "Updated Artist",
		Description:   "Updated Description",
		FormationYear: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	mockRepo := &mocks.MockArtistRepo{
		GetByIDFunc: func(ctx context.Context, id uint) (*model.Artist, error) {
			return originalArtist, nil
		},
		UpdateFunc: func(ctx context.Context, entity *model.Artist) (*model.Artist, error) {
			return updatedArtist, nil
		},
	}
	service := NewArtistService(mockRepo, logger)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	req := request.UpdateArtistRequest{
		ArtistName:    "Updated Artist",
		Description:   "Updated Description",
		FormationYear: "2021-01-01",
	}
	artist, err := service.UpdateArtist(ctx, 1, req)
	assert.NoError(t, err)
	assert.NotNil(t, artist)
	assert.Equal(t, updatedArtist, artist)
}