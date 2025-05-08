package service

import (
	"context"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"music-lib/internal/dto/response"
	"music-lib/internal/model"
	"music-lib/internal/service/mocks"
)

// TestSearchService_Search_Success проверяет успешный поиск по всем типам
func TestSearchService_Search_Success(t *testing.T) {
	// Создаем мок-репозитории
	mockArtistRepo := &mocks.MockArtistRepo{
		SearchFunc: func(ctx context.Context, query string, limit, offset int) ([]model.Artist, int64, error) {
			return []model.Artist{
				{ID: 1, Name: "Artist 1"},
				{ID: 2, Name: "Artist 2"},
			}, 2, nil
		},
	}
	mockAlbumRepo := &mocks.MockAlbumRepo{
		SearchFunc: func(ctx context.Context, query string, limit, offset int) ([]model.Album, int64, error) {
			return []model.Album{
				{ID: 1, Title: "Album 1"},
			}, 1, nil
		},
	}
	mockSongRepo := &mocks.MockSongRepo{
		SearchFunc: func(ctx context.Context, query string, limit, offset int) ([]model.Song, int64, error) {
			return []model.Song{
				{ID: 1, Title: "Song 1"},
			}, 1, nil
		},
	}

	// Создаем сервис
	service := NewSearchService(mockSongRepo, mockAlbumRepo, mockArtistRepo)

	// Создаем тестовый контекст Gin
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Выполняем поиск
	result := service.Search(ctx, []string{"artist", "album", "song"}, "test", 10, 0)

	// Проверяем результат
	searchResult, ok := result.(response.SearchResult)
	assert.True(t, ok, "Результат должен быть типа SearchResult")
	assert.Len(t, searchResult, 3, "Должно быть 3 типа в результате")

	// Проверяем данные для artist
	artistResp, exists := searchResult["artist"]
	assert.True(t, exists, "Тип 'artist' должен присутствовать")
	artistData, ok := artistResp.Data.([]response.ArtistDTO)
	assert.True(t, ok, "Данные для 'artist' должны быть []ArtistDTO")
	assert.Len(t, artistData, 2, "Должно быть 2 артиста")
	assert.Equal(t, uint(1), artistData[0].ID)
	assert.Equal(t, "Artist 1", artistData[0].Name)
	assert.Equal(t, 10, artistResp.Pagination.Limit)
	assert.Equal(t, 0, artistResp.Pagination.Offset)
	assert.Equal(t, int64(2), artistResp.Pagination.Total)

	// Проверяем данные для album
	albumResp, exists := searchResult["album"]
	assert.True(t, exists, "Тип 'album' должен присутствовать")
	albumData, ok := albumResp.Data.([]response.AlbumDTO)
	assert.True(t, ok, "Данные для 'album' должны быть []AlbumDTO")
	assert.Len(t, albumData, 1, "Должен быть 1 альбом")
	assert.Equal(t, uint(1), albumData[0].ID)
	assert.Equal(t, "Album 1", albumData[0].Title)

	// Проверяем данные для song
	songResp, exists := searchResult["song"]
	assert.True(t, exists, "Тип 'song' должен присутствовать")
	songData, ok := songResp.Data.([]response.SongDTO)
	assert.True(t, ok, "Данные для 'song' должны быть []SongDTO")
	assert.Len(t, songData, 1, "Должна быть 1 песня")
	assert.Equal(t, uint(1), songData[0].ID)
	assert.Equal(t, "Song 1", songData[0].Title)
}

// TestSearchService_Search_EmptyTypes проверяет поведение при пустом списке типов
func TestSearchService_Search_EmptyTypes(t *testing.T) {
	// Создаем мок-репозитории
	service := NewSearchService(&mocks.MockSongRepo{}, &mocks.MockAlbumRepo{}, &mocks.MockArtistRepo{})

	// Создаем тестовый контекст Gin
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Выполняем поиск с пустым списком типов
	result := service.Search(ctx, []string{}, "test", 10, 0)

	// Проверяем результат
	searchResult, ok := result.(response.SearchResult)
	assert.True(t, ok, "Результат должен быть типа SearchResult")
	assert.Len(t, searchResult, 0, "Результат должен быть пустым")
}

// TestSearchService_Search_UnknownType проверяет обработку неизвестного типа
func TestSearchService_Search_UnknownType(t *testing.T) {
	// Создаем мок-репозиторий для artist
	mockArtistRepo := &mocks.MockArtistRepo{
		SearchFunc: func(ctx context.Context, query string, limit, offset int) ([]model.Artist, int64, error) {
			return []model.Artist{{ID: 1, Name: "Artist 1"}}, 1, nil
		},
	}
	service := NewSearchService(nil, nil, mockArtistRepo)

	// Создаем тестовый контекст Gin
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Выполняем поиск с известным и неизвестным типом
	result := service.Search(ctx, []string{"artist", "unknown"}, "test", 10, 0)

	// Проверяем результат
	searchResult, ok := result.(response.SearchResult)
	assert.True(t, ok, "Результат должен быть типа SearchResult")
	assert.Len(t, searchResult, 2, "Должно быть 2 типа в результате")

	// Проверяем artist
	artistResp, exists := searchResult["artist"]
	assert.True(t, exists, "Тип 'artist' должен присутствовать")
	artistData, ok := artistResp.Data.([]response.ArtistDTO)
	assert.True(t, ok, "Данные для 'artist' должны быть []ArtistDTO")
	assert.Len(t, artistData, 1)

	// Проверяем unknown
	unknownResp, exists := searchResult["unknown"]
	assert.True(t, exists, "Тип 'unknown' должен присутствовать")
	errorResp, ok := unknownResp.Data.(response.SearchErrorResponse)
	assert.True(t, ok, "Данные для 'unknown' должны быть SearchErrorResponse")
	assert.Equal(t, "unknown search type", errorResp.Error.Error())
}

// TestSearchService_Search_RepoError проверяет обработку ошибки от репозитория
func TestSearchService_Search_RepoError(t *testing.T) {
	// Создаем мок-репозитории с ошибкой для artist
	mockArtistRepo := &mocks.MockArtistRepo{
		SearchFunc: func(ctx context.Context, query string, limit, offset int) ([]model.Artist, int64, error) {
			return nil, 0, errors.New("search error")
		},
	}
	mockAlbumRepo := &mocks.MockAlbumRepo{
		SearchFunc: func(ctx context.Context, query string, limit, offset int) ([]model.Album, int64, error) {
			return []model.Album{{ID: 1, Title: "Album 1"}}, 1, nil
		},
	}
	mockSongRepo := &mocks.MockSongRepo{
		SearchFunc: func(ctx context.Context, query string, limit, offset int) ([]model.Song, int64, error) {
			return []model.Song{{ID: 1, Title: "Song 1"}}, 1, nil
		},
	}

	service := NewSearchService(mockSongRepo, mockAlbumRepo, mockArtistRepo)

	// Создаем тестовый контекст Gin
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Выполняем поиск
	result := service.Search(ctx, []string{"artist", "album", "song"}, "test", 10, 0)

	// Проверяем результат
	searchResult, ok := result.(response.SearchResult)
	assert.True(t, ok, "Результат должен быть типа SearchResult")
	assert.Len(t, searchResult, 2, "Должно быть 2 типа (artist исключен из-за ошибки)")
	assert.NotContains(t, searchResult, "artist", "Тип 'artist' не должен присутствовать из-за ошибки")
	assert.Contains(t, searchResult, "album", "Тип 'album' должен присутствовать")
	assert.Contains(t, searchResult, "song", "Тип 'song' должен присутствовать")
}

// TestConvertToDTO проверяет функцию преобразования в DTO
func TestConvertToDTO(t *testing.T) {
	// Успешное преобразование для artist
	artists := []model.Artist{{ID: 1, Name: "Artist 1"}}
	dto := convertToDTO("artist", artists)
	artistDTOs, ok := dto.([]response.ArtistDTO)
	assert.True(t, ok, "Данные должны быть []ArtistDTO")
	assert.Len(t, artistDTOs, 1)
	assert.Equal(t, uint(1), artistDTOs[0].ID)
	assert.Equal(t, "Artist 1", artistDTOs[0].Name)

	// Некорректный тип данных
	dto = convertToDTO("artist", "not an artist slice")
	errorResp, ok := dto.(response.SearchErrorResponse)
	assert.True(t, ok, "Данные должны быть SearchErrorResponse")
	assert.Equal(t, "can't convert to artist model", errorResp.Error.Error())

	// Неизвестный тип
	dto = convertToDTO("unknown", nil)
	errorResp, ok = dto.(response.SearchErrorResponse)
	assert.True(t, ok, "Данные должны быть SearchErrorResponse")
	assert.Equal(t, "unknown search type", errorResp.Error.Error())
}