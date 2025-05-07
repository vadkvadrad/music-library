package v1

import (
	"music-lib/internal/dto/request"
	"music-lib/internal/middleware"
	"music-lib/pkg/er"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initProfileRoutes(api *gin.RouterGroup) {
	profile := api.Group("/profile")
	profile.Use(middleware.AuthMiddleware(h.config))
	{
		profile.POST("", h.NewProfile())
		profile.GET("", h.GetProfile())
	}
}


func (h *Handler) NewProfile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body := request.NewProfileRequest{}

		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.Error(er.ValidationError{Message: err.Error()})
			return
		}

		user, ok := middleware.GetUserData(ctx)
		if !ok {
			ctx.Error(er.ErrNotAuthorized)
			return
		}

		err := h.services.Profile.NewProfile(ctx, body, user.Id)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, nil)
	}
}


func (h *Handler) GetProfile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, ok := middleware.GetUserData(ctx)
		if !ok {
			ctx.Error(er.ErrNotAuthorized)
			return
		}

		profile, err := h.services.Profile.GetProfile(ctx, user.Id)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusOK, profile)
	}
}