package v1

import (
	"music-lib/internal/dto/request"
	"music-lib/internal/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initGenreRoutes(api *gin.RouterGroup) {
	genre := api.Group("/genre")
	genre.Use(middleware.AuthMiddleware(h.config))
	{
		genre.POST("", h.NewGenre())
	}
}


func (h *Handler) NewGenre() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body request.NewGenreRequest
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.Error(err)
			return
		}

		err := h.services.Genre.NewGenre(ctx, body.Name)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, nil)
	}
}