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
	Auth   *AuthService
	Album  *AlbumService
	Artist *ArtistService
	Search *SearchService
	Profile *ProfileService
}

func NewServices(deps *Deps) *Services {
	return &Services{
		Auth:   NewAuthService(deps.Repositories.User, deps.Event),
		Album:  NewAlbumService(deps.Repositories.Album, deps.Repositories.Artist),
		Artist: NewArtistService(deps.Repositories.Artist),
		Search: NewSearchService(deps.Repositories.Song, deps.Repositories.Album, deps.Repositories.Artist),
		Profile: NewProfileService(deps.Repositories.Profile),
	}
}
