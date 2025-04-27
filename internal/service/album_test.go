package service

import (
	"context"
	"errors"
	"fmt"
	"music-lib/internal/dto/request"
	"music-lib/internal/model"
	"music-lib/internal/service/mocks"
	"music-lib/pkg/er"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func TestAlbumService_NewAlbum(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    ctx, _ := gin.CreateTestContext(w)

    now := time.Now()
    testArtist := &model.Artist{
        ID: 1,
        Albums: []model.Album{
            {Title: "Existing Album"},
        },
    }

    tests := []struct {
        name        string
        setupMocks  func() (*mocks.MockArtistRepo, *mocks.MockAlbumRepo)
        request     request.NewAlbumRequest
        userID      uint
        expectedErr error
    }{
        {
            name: "Success - album created",
            setupMocks: func() (*mocks.MockArtistRepo, *mocks.MockAlbumRepo) {
                artistRepo := &mocks.MockArtistRepo{
                    GetByUserIDFunc: func(_ context.Context, uid uint) (*model.Artist, error) {
                        return &model.Artist{ID: 1, Albums: []model.Album{}}, nil
                    },
                }
                albumRepo := &mocks.MockAlbumRepo{
                    CreateFunc: func(_ context.Context, album *model.Album) error {
                        if album.ArtistID != 1 || album.Title != "New Album" {
                            return fmt.Errorf("invalid album data")
                        }
                        return nil
                    },
                }
                return artistRepo, albumRepo
            },
            request: request.NewAlbumRequest{
                Title:       "New Album",
                ReleaseDate: now.Format("2006-01-02"),
            },
            userID:      1,
            expectedErr: nil,
        },
        {
            name: "Artist not found",
            setupMocks: func() (*mocks.MockArtistRepo, *mocks.MockAlbumRepo) {
                artistRepo := &mocks.MockArtistRepo{
                    GetByUserIDFunc: func(_ context.Context, uid uint) (*model.Artist, error) {
                        return nil, gorm.ErrRecordNotFound
                    },
                }
                return artistRepo, &mocks.MockAlbumRepo{}
            },
            request:     request.NewAlbumRequest{ReleaseDate: now.Format("2006-01-02")},
            userID:      1,
            expectedErr: er.ErrArtistNotExists,
        },
        {
            name: "Invalid date format",
            setupMocks: func() (*mocks.MockArtistRepo, *mocks.MockAlbumRepo) {
                artistRepo := &mocks.MockArtistRepo{
                    GetByUserIDFunc: func(_ context.Context, uid uint) (*model.Artist, error) {
                        return testArtist, nil
                    },
                }
                return artistRepo, &mocks.MockAlbumRepo{}
            },
            request: request.NewAlbumRequest{
                Title:       "New Album",
                ReleaseDate: "invalid-date",
            },
            userID:      1,
            expectedErr: er.ErrDateFormat,
        },
        {
            name: "Album already exists",
            setupMocks: func() (*mocks.MockArtistRepo, *mocks.MockAlbumRepo) {
                artistRepo := &mocks.MockArtistRepo{
                    GetByUserIDFunc: func(_ context.Context, uid uint) (*model.Artist, error) {
                        return testArtist, nil
                    },
                }
                return artistRepo, &mocks.MockAlbumRepo{}
            },
            request: request.NewAlbumRequest{
                Title:       "Existing Album",
                ReleaseDate: now.Format("2006-01-02"),
            },
            userID:      1,
            expectedErr: er.ErrAlbumExists,
        },
        {
            name: "Error creating album",
            setupMocks: func() (*mocks.MockArtistRepo, *mocks.MockAlbumRepo) {
                artistRepo := &mocks.MockArtistRepo{
                    GetByUserIDFunc: func(_ context.Context, uid uint) (*model.Artist, error) {
                        return &model.Artist{ID: 1}, nil
                    },
                }
                albumRepo := &mocks.MockAlbumRepo{
                    CreateFunc: func(_ context.Context, _ *model.Album) error {
                        return fmt.Errorf("database error")
                    },
                }
                return artistRepo, albumRepo
            },
            request: request.NewAlbumRequest{
                Title:       "New Album",
                ReleaseDate: now.Format("2006-01-02"),
            },
            userID:      1,
            expectedErr: &er.InternalError{},
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            artistRepo, albumRepo := tt.setupMocks()
            service := NewAlbumService(albumRepo, artistRepo)

            err := service.NewAlbum(ctx, tt.request, tt.userID)

            switch {
            case tt.expectedErr == nil && err != nil:
                t.Fatalf("Unexpected error: %v", err)
            
            case tt.expectedErr != nil && err == nil:
                t.Fatal("Expected error but got nil")
            
            case tt.expectedErr != nil && !errors.Is(err, tt.expectedErr):
                if _, ok := tt.expectedErr.(*er.InternalError); ok {
                    if !strings.Contains(err.Error(), "database error") {
                        t.Errorf("Expected InternalError, got: %T", err)
                    }
                } else {
                    t.Errorf("Expected error: %v, got: %v", tt.expectedErr, err)
                }
            }
        })
    }
}

func TestAlbumService_GetAlbum(t *testing.T) {
    gin.SetMode(gin.TestMode)
    w := httptest.NewRecorder()
    ctx, _ := gin.CreateTestContext(w)

    testAlbum := &model.Album{
        ID:     1,
        Title:  "Test Album",
        Songs:  []model.Song{{ID: 1, Title: "Song1"}},
    }

    tests := []struct {
        name         string
        setupMock    func() *mocks.MockAlbumRepo
        strID        string
        expected     *model.Album
        expectedErr  error
        errorMessage string
    }{
        {
            name: "Success - valid album",
            setupMock: func() *mocks.MockAlbumRepo {
                return &mocks.MockAlbumRepo{
                    GetWithSongsFunc: func(_ context.Context, id uint) (*model.Album, error) {
                        if id == 1 {
                            return testAlbum, nil
                        }
                        return nil, fmt.Errorf("unexpected id")
                    },
                }
            },
            strID:       "1",
            expected:    testAlbum,
            expectedErr: nil,
        },
        {
            name: "Invalid ID - non-numeric",
            setupMock: func() *mocks.MockAlbumRepo {
                return &mocks.MockAlbumRepo{}
            },
            strID:        "invalid",
            expected:     nil,
            expectedErr:  &er.ValidationError{},
            errorMessage: "invalid syntax",
        },
        {
            name: "Album not found",
            setupMock: func() *mocks.MockAlbumRepo {
                return &mocks.MockAlbumRepo{
                    GetWithSongsFunc: func(_ context.Context, id uint) (*model.Album, error) {
                        return nil, gorm.ErrRecordNotFound
                    },
                }
            },
            strID:        "999",
            expected:     nil,
            expectedErr:  er.ErrAlbumNotExists,
        },
        {
            name: "Database error",
            setupMock: func() *mocks.MockAlbumRepo {
                return &mocks.MockAlbumRepo{
                    GetWithSongsFunc: func(_ context.Context, id uint) (*model.Album, error) {
                        return nil, fmt.Errorf("database connection failed")
                    },
                }
            },
            strID:        "1",
            expected:     nil,
            expectedErr:  er.InternalError{},
            errorMessage: "database connection failed",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            albumRepo := tt.setupMock()
            service := NewAlbumService(albumRepo, nil) // ArtistRepo не используется

            result, err := service.GetAlbum(ctx, tt.strID)

            // Проверка ошибок
            switch {
            case tt.expectedErr == nil && err != nil:
                t.Fatalf("Unexpected error: %v", err)
            
            case tt.expectedErr != nil && err == nil:
                t.Fatal("Expected error but got nil")
            
            case tt.expectedErr != nil:
                if !errors.Is(err, tt.expectedErr) {
                    // Для проверки кастомных типов ошибок
                    switch tt.expectedErr.(type) {
                    case *er.ValidationError:
                        if verr, ok := err.(*er.ValidationError); ok {
                            if !strings.Contains(verr.Message, tt.errorMessage) {
                                t.Errorf("Expected message containing '%s', got: '%s'", tt.errorMessage, verr.Message)
                            }
                        } else {
                            t.Errorf("Expected ValidationError, got: %T", err)
                        }
                    
                    case er.InternalError:
                        if ierr, ok := err.(er.InternalError); ok {
                            if !strings.Contains(ierr.Message, tt.errorMessage) {
                                t.Errorf("Expected message containing '%s', got: '%s'", tt.errorMessage, ierr.Message)
                            }
                        } else {
                            t.Errorf("Expected InternalError, got: %T", err)
                        }
                    
                    default:
                        t.Errorf("Expected error: %T, got: %T", tt.expectedErr, err)
                    }
                }
            }

            // Проверка результата
            if tt.expected != nil {
                if result.ID != tt.expected.ID || result.Title != tt.expected.Title {
                    t.Errorf("Expected album: %+v, got: %+v", tt.expected, result)
                }
                
                if len(result.Songs) != len(tt.expected.Songs) {
                    t.Errorf("Expected %d songs, got %d", len(tt.expected.Songs), len(result.Songs))
                }
            }
        })
    }
}