package service

import (
	"music-lib/internal/repository"
	"music-lib/pkg/event"

	"go.uber.org/zap"
)

type Deps struct {
	Event        *event.EventBus
	Repositories *repository.Repositories
	Logger       *zap.SugaredLogger
}

type Services struct {
	Auth       *AuthService
	Album      *AlbumService
	Artist     *ArtistService
	Song       *SongService
	Genre      *GenreService
	Search     *SearchService
	Profile    *ProfileService
	Permission *PermissionService
}

func NewServices(deps *Deps) *Services {
	return &Services{
		Auth:   NewAuthService(deps.Repositories.User, deps.Event),
		Album:  NewAlbumService(deps.Repositories.Album, deps.Repositories.Artist),
		Artist: NewArtistService(deps.Repositories.Artist, deps.Logger),
		Song: NewSongService(deps.Repositories.Song,
			deps.Repositories.Album,
			deps.Repositories.SongGenre,
			deps.Repositories.Genre,
			deps.Repositories.Lyrics,
			deps.Logger,
		),
		Genre:   NewGenreService(deps.Repositories.Genre, deps.Logger),
		Search:  NewSearchService(deps.Repositories.Song, deps.Repositories.Album, deps.Repositories.Artist),
		Profile: NewProfileService(deps.Repositories.Profile),
		Permission: NewPermissionService(deps.Repositories.Permission, deps.Logger),
	}
}
