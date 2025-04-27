package service

import (
	"context"
	"errors"
	"fmt"
	"music-lib/internal/model"
	"music-lib/internal/service/mocks"
	"music-lib/pkg/er"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestArtistService_GetArtist(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	testArtist := &model.Artist{
		ID:   1,
		Name: "Test Artist",
		Albums: []model.Album{
			{ID: 1, Title: "Album 1"},
		},
	}

	tests := []struct {
		name        string
		setupMock   func() *mocks.MockArtistRepo
		strID       string
		expected    *model.Artist
		expectedErr error
	}{
		{
			name: "Success",
			setupMock: func() *mocks.MockArtistRepo {
				return &mocks.MockArtistRepo{
					GetWithAlbumsFunc: func(_ context.Context, id uint) (*model.Artist, error) {
						if id == 1 {
							return testArtist, nil
						}
						return nil, fmt.Errorf("invalid id")
					},
				}
			},
			strID:       "1",
			expected:    testArtist,
			expectedErr: nil,
		},
		{
			name: "Invalid ID",
			setupMock: func() *mocks.MockArtistRepo {
				return &mocks.MockArtistRepo{}
			},
			strID:       "invalid",
			expectedErr: &er.ValidationError{},
		},
		{
			name: "Artist not found",
			setupMock: func() *mocks.MockArtistRepo {
				return &mocks.MockArtistRepo{
					GetWithAlbumsFunc: func(_ context.Context, id uint) (*model.Artist, error) {
						return nil, gorm.ErrRecordNotFound
					},
				}
			},
			strID:       "999",
			expectedErr: er.ErrArtistNotExists,
		},
		{
			name: "Database error",
			setupMock: func() *mocks.MockArtistRepo {
				return &mocks.MockArtistRepo{
					GetWithAlbumsFunc: func(_ context.Context, id uint) (*model.Artist, error) {
						return nil, fmt.Errorf("connection failed")
					},
				}
			},
			strID:       "1",
			expectedErr: &er.InternalError{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := tt.setupMock()
			service := NewArtistService(repo)

			result, err := service.GetArtist(ctx, tt.strID)

			// Проверка ошибок
			switch {
			case tt.expectedErr == nil && err != nil:
				t.Fatalf("Unexpected error: %v", err)

			case tt.expectedErr != nil && err == nil:
				t.Fatal("Expected error but got nil")

			case tt.expectedErr != nil:
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("Expected error: %T, got: %T", tt.expectedErr, err)
				}
			}

			// Проверка результата
			if tt.expected != nil {
				if !reflect.DeepEqual(result, tt.expected) {
					t.Errorf("Expected: %+v, got: %+v", tt.expected, result)
				}
			}
		})
	}
}