package repository

import (
	"music-lib/internal/model"
	"music-lib/internal/repository/postgres"
	"music-lib/pkg/db"
)

const (
	EmailKey     = "email"
	PhoneKey     = "phone"
	SessionIdKey = "session_id"
)

type ISongRepository interface {
	Create(song *model.Song) (*model.Song, error)
	Update(song *model.Song) (*model.Song, error)
	Delete(id uint) error
	FindByID(id uint) (*model.Song, error)
	Find(song, group string) (*model.Song, error)
	FindByGroup(group string, limit, offset int) []model.Song
}

type IUserRepository interface {
	Create(user *model.User) (*model.User, error)
	Update(user *model.User) (*model.User, error)
	FindByKey(key, data string) (*model.User, error)
}

type Repositories struct {
	Song ISongRepository
	User IUserRepository
}

func NewPostgresRepositories(db *db.Db) *Repositories {
	return &Repositories{
		Song: postgres.NewSongRepository(db),
		User: postgres.NewUserRepository(db),
	}
}
