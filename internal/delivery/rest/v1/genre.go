package v1

import (
	"music-lib/internal/dto/request"
	"music-lib/internal/dto/response"
	"music-lib/internal/middleware"
	"music-lib/pkg/er"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initGenreRoutes(api *gin.RouterGroup) {
	genre := api.Group("/genre")
	{
		genre.GET("", h.GetAllGenres())
		genre.GET("/:id", h.GetGenre())
	}
	genre.Use(middleware.AuthMiddleware(h.config))
	{
		genre.POST("", h.NewGenre())
		genre.PATCH("/:id", h.UpdateGenre())
		genre.DELETE("", h.DeleteGenre())
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

func (h *Handler) UpdateGenre() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body request.UpdateGenreRequest
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.Error(err)
			return
		}

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			h.logger.Debugw("Wrong id parameter",
				"error got:", err.Error(),
			)
			ctx.Error(er.ValidationError{Message: err.Error()})
			return
		}

		err = h.services.Genre.UpdateGenre(ctx, uint(id), body.NewName)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, nil)
	}
}

func (h *Handler) GetAllGenres() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		genres, err := h.services.Genre.GetAllGenres(ctx)
		if err != nil {
			ctx.Error(err)
			return
		}

		genreDTOs := make([]response.GenreDTO, len(genres))
		for i, genre := range genres {
			genreDTOs[i] = response.GenreDTO{
				ID:   genre.ID,
				Name: genre.Name,
			}
		}

		ctx.JSON(http.StatusOK, genreDTOs)
	}
}

func (h *Handler) GetGenre() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.Error(&er.ValidationError{Message: err.Error()})
			return
		}

		genre, err := h.services.Genre.GetGenre(ctx, uint(id))
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, response.GenreDTO{
			ID:   genre.ID,
			Name: genre.Name,
		})
	}
}

func (h *Handler) DeleteGenre() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// TODO: implement delete genre
		ctx.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
	}
}
