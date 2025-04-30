package app

import (
	"fmt"
	"music-lib/internal/config"
	"music-lib/internal/delivery/rest"
	"music-lib/internal/infrastructure/email"
	"music-lib/internal/repository"
	"music-lib/internal/service"
	"music-lib/pkg/db"
	"music-lib/pkg/event"

	"go.uber.org/zap"
)

func RunV1() {
	// Config
	cfg, err := config.Load()
	if err != nil {
		panic("bad config params: " + err.Error())
	}

	// Logger 
	logger, _ := zap.NewDevelopment()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	// Database
	dbCfg := db.NewDbConfig(cfg.Db.Dsn)
	db := db.NewDb(dbCfg)

	// Event
	eventBus := event.NewEventBus()

	// Repositories
	repositories := repository.NewPostgresRepositories(db)

	// Services
	services := service.NewServices(&service.Deps{
		Event: eventBus,
		Repositories: repositories,
		Logger: sugar,
	})

	// Handlers
	handlers := rest.NewHandler(services, cfg, sugar)
	router := handlers.Init(cfg)

	sender, err := email.Load(cfg, eventBus)
	if err != nil {
		panic("sender unable " + err.Error())
	}
	
	go sender.Listen()

	// Router run
	fmt.Println("Server started on port ", cfg.App.Port)
	router.Run(":" + cfg.App.Port)
}