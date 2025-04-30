package v1

import (
	"music-lib/internal/config"
	"music-lib/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	services *service.Services
	config   *config.Config
	logger   *zap.SugaredLogger
}

func NewHandler(services *service.Services, conf *config.Config, zapLogger *zap.SugaredLogger) *Handler {
	return &Handler{
		services: services,
		config:   conf,
		logger: zapLogger,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initAuthRoutes(v1)
		h.initSearchRoutes(v1)
		h.initSongRoutes(v1)
		h.initProfileRoutes(v1)
		h.initArtistRoutes(v1)
		h.initAlbumRoutes(v1)
	}
}
