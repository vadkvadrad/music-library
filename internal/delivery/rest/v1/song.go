package v1

import (
	"music-lib/internal/dto/request"
	"music-lib/internal/dto/response"
	"music-lib/internal/middleware"
	"music-lib/internal/model"
	"music-lib/pkg/er"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initSongRoutes(api *gin.RouterGroup) {
	song := api.Group("/song")
	song.GET("/:id", h.GetSong())
	song.Use(middleware.AuthMiddleware(h.config))
	{
		song.POST("/:album_id", h.AddSong())
		song.PATCH("/:id", h.UpdateSong())
		song.GET("/:id/has-permission", h.HasSongPermission())
	}
}

func (h *Handler) AddSong() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		albumID, err := strconv.Atoi(ctx.Param("album_id"))
		if err != nil {
			h.logger.Debug("Invalid album ID format",
				"received", ctx.Param("album_id"),
				"error", err,
			)
			ctx.Error(&er.ValidationError{Message: err.Error()})
			return
		}

		var body request.NewSongRequest
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.Error(err)
			return
		}

		user, ok := middleware.GetUserData(ctx)
		if !ok {
			h.logger.Debug("User credentials not found",
				"error", er.ErrWrongUserCredentials.Error(),
			)
			ctx.Error(er.ErrNotAuthorized)
			return
		}

		h.logger.Infow("Attempting to get album for song creation",
			"user_id", user.Id,
			"album_id", albumID,
		)

		album, err := h.services.Album.GetArtistAlbum(ctx, user.Id, uint(albumID))
		if err != nil {
			h.logger.Errorw("Failed to get album",
				"user_id", user.Id,
				"album_id", albumID,
				"error", err.Error(),
			)
			ctx.Error(err)
			return
		}

		h.logger.Infow("Album found successfully",
			"album_id", album.ID,
			"album_title", album.Title,
			"artist_id", album.ArtistID,
		)

		h.logger.Infow("Adding new song",
			"song_title", body.Title,
			"album_id", album.ID,
			"artist_id", album.ArtistID,
		)

		song, err := h.services.Song.AddSong(ctx, album, body)
		if err != nil {
			ctx.Error(err)
			return
		}

		// Добавить разрешение для песен
		err = h.services.Permission.AddPermission(ctx, user.Id, song.ID, model.SongResource, model.EditPermission)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, response.AddSongResponse{
			AlbumID:   album.ID,
			AlbumName: album.Title,
			ArtistID:  album.ArtistID,
		})
	}
}

func (h *Handler) GetSong() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			h.logger.Debug("Invalid ID format",
				"received", ctx.Param("id"),
				"error", err,
			)
			ctx.Error(&er.ValidationError{Message: err.Error()})
			return
		}

		song, lyrics, genres, err := h.services.Song.GetSong(ctx, uint(id))
		if err != nil {
			ctx.Error(err)
			return
		}

		var albumID uint
		if song.AlbumID != nil {
			albumID = *song.AlbumID
		}

		// Преобразуем лирику в DTO
		var lyricsDTO response.LyricsDTO
		if lyrics != nil && len(lyrics.Couplets) > 0 {
			coupletsDTO := make([]response.CoupletDTO, len(lyrics.Couplets))
			for i, couplet := range lyrics.Couplets {
				coupletsDTO[i] = response.CoupletDTO{
					Number:  couplet.Number,
					Couplet: couplet.Text,
				}
			}
			lyricsDTO.Couplets = coupletsDTO
		}

		// Преобразуем жанры в DTO
		genreDTOs := make([]response.GenreDTO, len(genres))
		for i, genre := range genres {
			genreDTOs[i] = response.GenreDTO{
				ID:   genre.ID,
				Name: genre.Name,
			}
		}

		ctx.JSON(http.StatusOK, response.SongDTO{
			ID:       song.ID,
			Title:    song.Title,
			AlbumID:  albumID,
			Duration: song.Duration,
			FilePath: song.FilePath,
			Lyrics:   lyricsDTO,
			Genres:   genreDTOs,
		})
	}
}

func (h *Handler) UpdateSong() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			h.logger.Debug("Invalid ID format",
				"received", ctx.Param("id"),
				"error", err,
			)
			ctx.Error(&er.ValidationError{Message: err.Error()})
			return
		}

		var body request.UpdateSongRequest
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.Error(err)
			return
		}

		user, ok := middleware.GetUserData(ctx)
		if !ok {
			ctx.Error(er.ErrNotAuthorized)
			return
		}

		// Проверяем права на редактирование песни
		if !h.services.Permission.HasPermission(user.Id, uint(id), model.SongResource, model.EditPermission) {
			ctx.Error(er.ErrWrongUserCredentials)
			return
		}

		_, err = h.services.Song.UpdateSong(ctx, uint(id), body)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, nil)
	}
}

func (h *Handler) HasSongPermission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.Error(&er.ValidationError{Message: err.Error()})
			return
		}

		user, ok := middleware.GetUserData(ctx)
		if !ok {
			ctx.Error(er.ErrNotAuthorized)
			return
		}

		hasPermission := h.services.Permission.HasPermission(user.Id, uint(id), model.SongResource, model.EditPermission)

		ctx.JSON(http.StatusOK, gin.H{
			"has_permission": hasPermission,
		})
	}
}
