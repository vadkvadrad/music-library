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

func (h *Handler) initArtistRoutes(api *gin.RouterGroup) {
	artist := api.Group("/artist")
	{
		artist.GET("/:id", h.GetArtist())
	}
	artist.Use(middleware.AuthMiddleware(h.config))
	{
		artist.POST("", h.NewArtist())
		artist.PATCH("/:id", h.UpdateArtist())
		artist.GET("/:id/has-permission", h.HasArtistPermission())
	}
}

func (h *Handler) NewArtist() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body request.NewArtistRequest

		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.Error(err)
			return
		}

		user, ok := middleware.GetUserData(ctx)
		if !ok {
			ctx.Error(er.ErrNotAuthorized)
			return
		}

		artist, err := h.services.Artist.NewArtist(ctx, body, user.Id)
		if err != nil {
			ctx.Error(err)
			return
		}

		err = h.services.Permission.AddPermission(ctx, user.Id, artist.ID, model.ArtistResource, model.EditPermission)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, response.ArtistDTO{
			ID:            artist.ID,
			Name:          artist.Name,
			Description:   artist.Description,
			FormationYear: artist.FormationYear,
		})
	}
}

func (h *Handler) GetArtist() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		strID := ctx.Param("id")

		artist, err := h.services.Artist.GetArtist(ctx, strID)
		if err != nil {
			ctx.Error(err)
			return
		}

		// Загружаем альбомы артиста
		albums, err := h.services.Artist.GetAlbums(ctx, artist.ID)
		if err != nil {
			// Если ошибка при загрузке альбомов, возвращаем артиста без альбомов
			albums = []model.Album{}
		}

		// Преобразуем альбомы в DTO
		albumDTOs := make([]response.AlbumDTO, len(albums))
		for i, album := range albums {
			albumDTOs[i] = response.AlbumDTO{
				ID:          album.ID,
				Title:       album.Title,
				ReleaseDate: album.ReleaseDate,
				CoverArtURL: album.CoverArtURL,
			}
		}

		ctx.JSON(http.StatusOK, response.ArtistDTO{
			ID:            artist.ID,
			Name:          artist.Name,
			Description:   artist.Description,
			FormationYear: artist.FormationYear,
			Albums:        albumDTOs,
		})
	}
}

func (h *Handler) UpdateArtist() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		strID := ctx.Param("id")
		id, err := strconv.Atoi(strID)
		if err != nil {
			ctx.Error(&er.ValidationError{Message: err.Error()})
			return
		}

		var body request.UpdateArtistRequest
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.Error(err)
			return
		}

		h.logger.Infow("Updating artist",
			"id", id,
			"name to update", body.ArtistName,
			"description to update", body.Description,
			"formation year to update", body.FormationYear,
		)

		user, ok := middleware.GetUserData(ctx)
		if !ok {
			ctx.Error(er.ErrNotAuthorized)
			return
		}

		if !h.services.Permission.HasPermission(user.Id, uint(id), model.ArtistResource, model.EditPermission) {
			ctx.Error(er.ErrWrongUserCredentials)
			return
		}

		artist, err := h.services.Artist.UpdateArtist(ctx, uint(id), body)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, response.ArtistDTO{
			ID:            artist.ID,
			Name:          artist.Name,
			Description:   artist.Description,
			FormationYear: artist.FormationYear,
		})
	}
}

func (h *Handler) HasArtistPermission() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		strID := ctx.Param("id")
		id, err := strconv.Atoi(strID)
		if err != nil {
			ctx.Error(&er.ValidationError{Message: err.Error()})
			return
		}

		user, ok := middleware.GetUserData(ctx)
		if !ok {
			ctx.Error(er.ErrNotAuthorized)
			return
		}

		hasPermission := h.services.Permission.HasPermission(user.Id, uint(id), model.ArtistResource, model.EditPermission)

		ctx.JSON(http.StatusOK, gin.H{
			"has_permission": hasPermission,
		})
	}
}
