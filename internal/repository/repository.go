package repository

import (
	"context"
	"music-lib/internal/model"
	"music-lib/internal/repository/postgres"
	"music-lib/pkg/db"
)

const (
	EmailKey     = "email"
	PhoneKey     = "phone"
	SessionIdKey = "session_id"
)

// Базовый интерфейс для всех репозиториев
type Repository[T any] interface {
	Create(ctx context.Context, entity *T) (*T, error)
	Update(ctx context.Context, entity *T) (*T, error)
	Delete(ctx context.Context, id uint) error
}

type Searchable[T any] interface {
	Search(ctx context.Context, query string, limit, offset int) ([]T, int64, error)
}

// Репозиторий артистов
type IArtistRepository interface {
	Repository[model.Artist]
	Searchable[model.Artist]

	GetByID(ctx context.Context, id uint) (*model.Artist, error)
	GetByUserID(ctx context.Context, userID uint) (*model.Artist, error)
	GetWithAlbums(ctx context.Context, id uint) (*model.Artist, error)
	GetArtistAlbumByUserID(ctx context.Context, userID uint, albumID uint) (*model.Album, int, error)
	IsExists(ctx context.Context, name string) bool
}

// Репозиторий альбомов
type IAlbumRepository interface {
	Repository[model.Album]
	Searchable[model.Album]

	GetByID(ctx context.Context, id uint) (*model.Album, error)
	GetWithSongs(ctx context.Context, id uint) (*model.Album, error)
}

// Репозиторий песен
type ISongRepository interface {
	Repository[model.Song]
	Searchable[model.Song]

	ExistsInAlbum(ctx context.Context, albumID uint, songName string) bool
	GetByID(ctx context.Context, id uint) (*model.Song, error)
	GetByArtistID(ctx context.Context, artistID uint, sort string, limit, offset int) ([]model.Song, int64, error)
	GetByAlbumID(ctx context.Context, albumID uint, sort string, limit, offset int) ([]model.Song, int64, error)
	GetFullInfo(ctx context.Context, id uint) (*model.Song, *model.Artist, *model.Album, error)
}

// Репозиторий текстов песен
type ILyricsRepository interface {
	GetBySongID(ctx context.Context, songID uint) (*model.Lyrics, error)
	Upsert(ctx context.Context, lyrics *model.Lyrics) error
	DeleteBySongID(ctx context.Context, songID uint) error
}

type IUserRepository interface {
	Create(user *model.User) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	FindByKey(key, data string) (*model.User, error)
}

type IProfileRepository interface {
	Repository[model.Profile]

	GetByUserID(ctx context.Context, userID uint) (*model.Profile, error)
}

type IGenreRepository interface {
	Repository[model.Genre]

	GetById(ctx context.Context, id uint) (*model.Genre, error)
	GetByIds(ctx context.Context, ids []uint) ([]model.Genre, error)
	IsExists(ctx context.Context, name string) bool
}

type ISongGenreRepository interface {
	Repository[model.SongGenre]
}

type IPermissionRepository interface {
	Repository[model.ResourcePermission]

	HasPermission(userID, resourceID uint, resourceType model.Resource, permission model.Permission) bool
}

type Repositories struct {
	// User
	User IUserRepository
	// Music
	Song      ISongRepository
	Album     IAlbumRepository
	Artist    IArtistRepository
	Lyrics    ILyricsRepository
	Genre     IGenreRepository
	SongGenre ISongGenreRepository
	// Profile
	Profile IProfileRepository
	// Permission
	Permission IPermissionRepository
}

func NewPostgresRepositories(db *db.Db) *Repositories {
	return &Repositories{
		// User
		User: postgres.NewUserRepository(db),
		// Music
		Artist:    postgres.NewArtistRepository(db),
		Album:     postgres.NewAlbumRepository(db),
		Song:      postgres.NewSongRepository(db),
		Genre:     postgres.NewGenreRepository(db),
		SongGenre: postgres.NewSongGenreRepository(db),
		Lyrics:    postgres.NewLyricsRepository(db),
		//Profile
		Profile: postgres.NewProfileRepository(db),
		// Permission
		Permission: postgres.NewPermissionRepository(db),
	}
}
