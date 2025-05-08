package service

import (
	"context"
	"errors"
	"music-lib/internal/dto/request"
	"music-lib/internal/model"
	"music-lib/internal/service/mocks"
	"music-lib/pkg/er"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func TestAddSong_SongExists(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockSongRepo := &mocks.MockSongRepo{
		ExistsInAlbumFunc: func(ctx context.Context, albumID uint, songName string) bool {
			return true
		},
	}
	service := NewSongService(mockSongRepo, nil, nil, nil, nil, logger)
	album := &model.Album{ID: 1}
	req := request.NewSongRequest{Title: "Test Song"}

	song, err := service.AddSong(context.Background(), album, req)

	assert.Nil(t, song)
	assert.Equal(t, er.ErrSongExists, err)
}

func TestAddSong_CreateError(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockSongRepo := &mocks.MockSongRepo{
		ExistsInAlbumFunc: func(ctx context.Context, albumID uint, songName string) bool {
			return false
		},
		CreateFunc: func(ctx context.Context, entity *model.Song) (*model.Song, error) {
			return nil, errors.New("create error")
		},
	}
	service := NewSongService(mockSongRepo, nil, nil, nil, nil, logger)
	album := &model.Album{ID: 1}
	req := request.NewSongRequest{Title: "Test Song"}

	song, err := service.AddSong(context.Background(), album, req)

	assert.Nil(t, song)
	assert.IsType(t, &er.InternalError{}, err)
	assert.Equal(t, "create error", err.(*er.InternalError).Message)
}

func TestAddSong_Success(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockSongRepo := &mocks.MockSongRepo{
		ExistsInAlbumFunc: func(ctx context.Context, albumID uint, songName string) bool {
			return false
		},
		CreateFunc: func(ctx context.Context, entity *model.Song) (*model.Song, error) {
			return &model.Song{ID: 1}, nil
		},
	}
	mockGenreRepo := &mocks.MockGenreRepo{
		GetByIdsFunc: func(ctx context.Context, ids []uint) ([]model.Genre, error) {
			return []model.Genre{{ID: 1}}, nil
		},
	}
	mockSongGenreRepo := &mocks.MockSongGenreRepo{
		CreateFunc: func(ctx context.Context, entity *model.SongGenre) (*model.SongGenre, error) {
			return &model.SongGenre{}, nil
		},
	}
	mockLyricsRepo := &mocks.MockLyricsRepo{
		UpsertFunc: func(ctx context.Context, lyrics *model.Lyrics) error {
			return nil
		},
	}
	service := NewSongService(mockSongRepo, nil, mockSongGenreRepo, mockGenreRepo, mockLyricsRepo, logger)
	album := &model.Album{ID: 1, ArtistID: 1}
	req := request.NewSongRequest{
		Title:  "Test Song",
		Genres: []request.Genres{{GenreID: 1}},
		Lyrics: request.AddLyrics{Text: []request.Couplet{{Text: "Lyrics"}}},
	}

	song, err := service.AddSong(context.Background(), album, req)

	assert.NotNil(t, song)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), song.ID)
}

func TestAddGenres_GetGenresError(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockGenreRepo := &mocks.MockGenreRepo{
		GetByIdsFunc: func(ctx context.Context, ids []uint) ([]model.Genre, error) {
			return nil, errors.New("get genres error")
		},
	}
	service := NewSongService(nil, nil, nil, mockGenreRepo, nil, logger)

	err := service.addGenres(context.Background(), 1, []request.Genres{{GenreID: 1}})

	assert.IsType(t, &er.InternalError{}, err)
	assert.Equal(t, "get genres error", err.(*er.InternalError).Message)
}

func TestAddGenres_Success(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockGenreRepo := &mocks.MockGenreRepo{
		GetByIdsFunc: func(ctx context.Context, ids []uint) ([]model.Genre, error) {
			return []model.Genre{{ID: 1}}, nil
		},
	}
	mockSongGenreRepo := &mocks.MockSongGenreRepo{
		CreateFunc: func(ctx context.Context, entity *model.SongGenre) (*model.SongGenre, error) {
			return &model.SongGenre{}, nil
		},
	}
	service := NewSongService(nil, nil, mockSongGenreRepo, mockGenreRepo, nil, logger)

	err := service.addGenres(context.Background(), 1, []request.Genres{{GenreID: 1}})

	assert.NoError(t, err)
}

func TestAddLyrics_UpsertError(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockLyricsRepo := &mocks.MockLyricsRepo{
		UpsertFunc: func(ctx context.Context, lyrics *model.Lyrics) error {
			return errors.New("upsert error")
		},
	}
	service := NewSongService(nil, nil, nil, nil, mockLyricsRepo, logger)

	req := request.AddLyrics{Text: []request.Couplet{{Text: "Lyrics"}}}
	err := service.addLyrics(context.Background(), 1, req)

	assert.Equal(t, errors.New("upsert error"), err)
}

func TestAddLyrics_Success(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockLyricsRepo := &mocks.MockLyricsRepo{
		UpsertFunc: func(ctx context.Context, lyrics *model.Lyrics) error {
			return nil
		},
	}
	service := NewSongService(nil, nil, nil, nil, mockLyricsRepo, logger)

	req := request.AddLyrics{Text: []request.Couplet{{Text: "Lyrics"}}}
	err := service.addLyrics(context.Background(), 1, req)

	assert.NoError(t, err)
}

func TestGetSong_NotFound(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockSongRepo := &mocks.MockSongRepo{
		GetByIDFunc: func(ctx context.Context, id uint) (*model.Song, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	service := NewSongService(mockSongRepo, nil, nil, nil, nil, logger)

	song, err := service.GetSong(context.Background(), 1)

	assert.Nil(t, song)
	assert.Equal(t, er.ErrSongNotExists, err)
}

func TestGetSong_Success(t *testing.T) {
	logger := zap.NewNop().Sugar()
	mockSongRepo := &mocks.MockSongRepo{
		GetByIDFunc: func(ctx context.Context, id uint) (*model.Song, error) {
			return &model.Song{ID: 1}, nil
		},
	}
	service := NewSongService(mockSongRepo, nil, nil, nil, nil, logger)

	song, err := service.GetSong(context.Background(), 1)

	assert.NotNil(t, song)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), song.ID)
}