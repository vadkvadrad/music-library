package postgres

import (
	"database/sql"
	"errors"
	"fmt"
	"music-lib/internal/model"
	"music-lib/pkg/db"
	"time"
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
	query := `
		INSERT INTO users (name, email, password, role, session_id, code, is_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	err := repo.Db.QueryRow(
		query,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
		user.SessionId,
		user.Code,
		user.IsVerified,
		now,
		now,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) Update(user *model.User) (*model.User, error) {
	query := `
		UPDATE users 
		SET name = $1, email = $2, password = $3, role = $4, session_id = $5, 
		    code = $6, is_verified = $7, updated_at = $8
		WHERE id = $9 AND deleted_at IS NULL
		RETURNING updated_at
	`

	now := time.Now()
	err := repo.Db.QueryRow(
		query,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
		user.SessionId,
		user.Code,
		user.IsVerified,
		now,
		user.ID,
	).Scan(&user.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

func (repo *UserRepository) FindByKey(key, data string) (*model.User, error) {
	var user model.User
	var deletedAt sql.NullTime

	query := fmt.Sprintf(`
		SELECT id, created_at, updated_at, deleted_at, name, email, password, 
		       role, session_id, code, is_verified
		FROM users
		WHERE %s = $1 AND deleted_at IS NULL
	`, key)

	err := repo.Db.QueryRow(query, data).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&deletedAt,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.SessionId,
		&user.Code,
		&user.IsVerified,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	if deletedAt.Valid {
		user.DeletedAt = &deletedAt.Time
	}

	return &user, nil
}
