package v1

import (
	"music-lib/internal/dto/request"
	"music-lib/internal/dto/response"
	"music-lib/internal/middleware"
	"music-lib/pkg/er"
	"net/http"

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
			ctx.Error(er.ErrWrongUserCredentials)
			return
		}

		err := h.services.Artist.NewArtist(ctx, body, user.Id)
		if err != nil {
			ctx.Error(err)
		}
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

		var albums []response.AlbumDTO 
		for _, album := range artist.Albums {
			albums = append(albums, response.AlbumDTO{
				ID: album.ID,
				Title: album.Title,
				ReleaseDate: album.ReleaseDate,
				CoverArtURL: album.CoverArtURL,
				Songs: nil,
			})
		}

		ctx.JSON(http.StatusOK, response.ArtistDTO{
			ID: artist.ID,
			Name: artist.Name,
			Description: artist.Description,
			FormationYear: artist.FormationYear,
			Albums: albums,
		})
	}
}
