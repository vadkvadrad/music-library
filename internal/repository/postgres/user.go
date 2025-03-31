package postgres

import (
	"fmt"
	"music-lib/internal/model"
	"music-lib/pkg/db"

	"gorm.io/gorm/clause"
)

type UserRepository struct {
	Db *db.Db
}

func NewUserRepository(db *db.Db) *UserRepository {
	return &UserRepository{
		Db: db,
	}
}

func (repo *UserRepository) Create(user *model.User) (*model.User, error) {
	result := repo.Db.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (repo *UserRepository) Update(user *model.User) (*model.User, error) {
	result := repo.Db.Clauses(clause.Returning{}).Updates(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (repo *UserRepository) FindByKey(key, data string) (*model.User, error) {
	var user model.User
	query := fmt.Sprintf("%s = ?", key)
	result := repo.Db.First(&user, query, data)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
