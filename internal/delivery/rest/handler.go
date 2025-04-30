package rest

import (
	"music-lib/internal/config"
	"music-lib/internal/delivery/rest/v1"
	"music-lib/internal/middleware"
	"music-lib/internal/service"
	"music-lib/pkg/er"
	"net/http"

	docs "music-lib/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
	//"go.uber.org/zap"
)

type Handler struct {
	services *service.Services
	config   *config.Config
	logger   *zap.SugaredLogger
}

func NewHandler(services *service.Services, conf *config.Config, sugar *zap.SugaredLogger) *Handler {
	return &Handler{
		services: services,
		config:   conf,
		logger:   sugar,
	}
}

func (h *Handler) Init(conf *config.Config) *gin.Engine {
	//logger, _ := zap.NewProduction()

	errorHandler := er.NewErrorHandler(&er.ErrorHandlerConfig{
		LogErrors:    true,
		ShowInternal: true,
		AppName:      "music-lib",
	})

	// Init gin handler
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		middleware.Logger(),
		errorHandler.GinMiddleware(),
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
	handlers := v1.NewHandler(h.services, h.config, h.logger)
	api := router.Group("/api")
	{
		handlers.Init(api)
	}
}
