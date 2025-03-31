package v1

import (
	"music-lib/internal/config"
	"music-lib/internal/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services     *service.Services
	config *config.Config
}

func NewHandler(services *service.Services, conf *config.Config) *Handler {
	return &Handler{
		services:     services,
		config: conf,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initAuthRoutes(v1)
		h.initSongRoutes(v1)
	}
}