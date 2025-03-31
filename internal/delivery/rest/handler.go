package rest

import (
	"music-lib/internal/config"
	"music-lib/internal/delivery/rest/v1"
	"music-lib/internal/middleware"
	"music-lib/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
	docs "music-lib/docs"
)

type Handler struct {
	services *service.Services
	config   *config.Config
}

func NewHandler(services *service.Services, conf *config.Config) *Handler {
	return &Handler{
		services: services,
		config:   conf,
	}
}

func (h *Handler) Init(conf *config.Config) *gin.Engine {
	// Init gin handler
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		middleware.Logger(),
	)

	docs.SwaggerInfo.BasePath = "/api/v1"


	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	h.initApi(router)

	return router
}

func (h *Handler) initApi(router *gin.Engine) {
	handlers := v1.NewHandler(h.services, h.config)
	api := router.Group("/api")
	{
		handlers.Init(api)
	}
}
