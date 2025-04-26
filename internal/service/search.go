package service

import (
	"fmt"
	"music-lib/internal/dto/response"
	"music-lib/internal/model"
	"music-lib/internal/repository"
	"sync"

	"github.com/gin-gonic/gin"
)

type SearchService struct {
	songRepo   repository.ISongRepository
	albumRepo  repository.IAlbumRepository
	artistRepo repository.IArtistRepository
}

func NewSearchService(
	songRepo repository.ISongRepository,
	albumRepo repository.IAlbumRepository,
	artistRepo repository.IArtistRepository,
) *SearchService {
	return &SearchService{
		songRepo:   songRepo,
		albumRepo:  albumRepo,
		artistRepo: artistRepo,
	}
}



func (s *SearchService) Search(
	c *gin.Context, 
	types []string,
	query string,
	limit int,
	offset int,
) any {
	var result response.SearchResult
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, entityType := range types {
		wg.Add(1)
		
		go func(t string) {
			defer wg.Done()
			
			var data interface{}
			var total int64
			var err error

			switch t {
			case "artist":
				data, total, err = s.artistRepo.Search(c, query, limit, offset)
			case "album":
				data, total, err = s.albumRepo.Search(c, query, limit, offset)
			case "song":
				data, total, err = s.songRepo.Search(c, query, limit, offset)
			}

			if err != nil {
				return
			}

			mu.Lock()
			defer mu.Unlock()
			
			result[t] = response.PaginatedResponse{
				Data:       convertToDTO(t, data),
				Pagination: response.Pagination{Limit: limit, Offset: offset, Total: total},
			}
		}(entityType)
	}

	wg.Wait()
	return result
}

func convertToDTO(t string, data any) any {
	switch t {
	case "artist":
		artist, ok := data.(model.Artist)
		if !ok {
			return response.SearchErrorResponse{
				Error: fmt.Errorf("can't convert to artist model"),
				Tip: "check the validity of Artist model",
			}
		}
		return response.ArtistDTO{
			ID: artist.ID,
			Name: artist.Name,
			Description: artist.Description,
			FormationYear: artist.FormationYear,
			Albums: nil,
		}
	case "album":
		album, ok := data.(model.Album)
		if !ok {
			return response.SearchErrorResponse{
				Error: fmt.Errorf("can't convert to album model"),
				Tip: "check the validity of Album model ",
			}
		}
		return response.AlbumDTO{
			ID: album.ID,
			Title: album.Title,
			ReleaseDate: album.ReleaseDate,
			CoverArtURL: album.CoverArtURL,
		}
	case "song":
		song, ok := data.(model.Song)
		if !ok {
			return response.SearchErrorResponse{
				Error: fmt.Errorf("can't convert to song model"),
				Tip: "check the validity of Song model ",
			}
		}
		return response.SongDTO{
			ID: song.ID,
			Title: song.Title,
			Duration: song.Duration,
			FilePath: song.FilePath,
		}
	default:
		return response.SearchErrorResponse{
			Error: fmt.Errorf("unknown search type"),
			Tip: "check type conversion",
		}
	}
}