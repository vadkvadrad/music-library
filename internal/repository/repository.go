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
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
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

type Repositories struct {
	// User
	User IUserRepository
	// Music
	Song   ISongRepository
	Album  IAlbumRepository
	Artist IArtistRepository
	Lyrics ILyricsRepository
	// Profile
	Profile IProfileRepository
}

func NewPostgresRepositories(db *db.Db) *Repositories {
	return &Repositories{
		// User
		User: postgres.NewUserRepository(db),
		// Music
		Artist: postgres.NewArtistRepository(db),
		Album:  postgres.NewAlbumRepository(db),
		//Profile
		Profile: postgres.NewProfileRepository(db),
	}
}
