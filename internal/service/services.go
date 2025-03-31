package service

import (
	"music-lib/internal/repository"
	"music-lib/pkg/event"
)

type Deps struct {
	Event        *event.EventBus
	Repositories *repository.Repositories
}

type Services struct {
	Auth *AuthService
	Song *SongService
}

func NewServices(deps *Deps) *Services {
	return &Services{
		Auth: NewAuthService(deps.Repositories.User, deps.Event),
		Song: NewSongService(deps.Repositories.Song),
	}
}
