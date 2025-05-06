package v1

import (
	"music-lib/internal/dto/request"
	"music-lib/internal/dto/response"
	"music-lib/internal/middleware"
	"music-lib/internal/model"
	"music-lib/pkg/er"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initAlbumRoutes(api *gin.RouterGroup) {
	album := api.Group("/album")
	{
		album.GET("/:id", h.GetAlbum())
	}
	album.Use(middleware.AuthMiddleware(h.config))
	{
		album.POST("", h.NewAlbum())
	}
}


func (h *Handler) NewAlbum() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body request.NewAlbumRequest

		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.Error(err)
			return
		}

		user, ok := middleware.GetUserData(ctx)
		if !ok {
			ctx.Error(er.ErrWrongUserCredentials)
			return
		}

		album, err := h.services.Album.NewAlbum(ctx, body, user.Id)
		if err != nil {
			ctx.Error(err)
			return
		}

		// Добавление разрешения на редактирование
		h.services.Permission.AddPermission(ctx, user.Id, album.ID, model.AlbumResource, model.EditPermission)

		ctx.JSON(http.StatusCreated, nil)
	}
}


func (h *Handler) GetAlbum() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		strID := ctx.Param("id")

		album, err := h.services.Album.GetAlbum(ctx, strID)
		if err != nil {
			ctx.Error(err)
			return
		}

		var songs []response.SongDTO
		for _, song := range album.Songs {
			songs = append(songs, response.SongDTO{
				ID: song.ID,
				Title: song.Title,
				Duration: song.Duration,
				FilePath: song.FilePath,
			})
		}
		ctx.JSON(http.StatusOK, response.AlbumDTO{
			ID: album.ID,
			Title: album.Title,
			ReleaseDate: album.ReleaseDate,
			CoverArtURL: album.CoverArtURL,
			Songs: songs,
		})
	}
}