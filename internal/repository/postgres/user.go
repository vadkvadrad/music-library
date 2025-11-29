package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"music-lib/internal/model"
	"music-lib/pkg/db"
)

type UserRepository struct {
	Db *db.Db
}

func NewUserRepository(db *db.Db) *UserRepository {
	return &UserRepository{
		Db: db,
	}
}

func (repo *UserRepository) Create(ctx context.Context, user *model.User) (*model.User, error) {
	query := `
		INSERT INTO users (name, email, password, role, session_id, code, is_verified, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	var id uint
	err := repo.Db.QueryRowContext(ctx, query,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
		user.SessionId,
		user.Code,
		user.IsVerified,
		user.CreatedAt,
		user.UpdatedAt,
		user.DeletedAt,
	).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	user.ID = id
	return user, nil
}

func (repo *UserRepository) Update(ctx context.Context, user *model.User) (*model.User, error) {
	query := `
		UPDATE users
		SET name = $1, email = $2, password = $3, role = $4, session_id = $5, code = $6, is_verified = $7, updated_at = $8
		WHERE id = $9
	`

	result, err := repo.Db.ExecContext(ctx, query,
		user.Name,
		user.Email,
		user.Password,
		user.Role,
		user.SessionId,
		user.Code,
		user.IsVerified,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("no user found with id %d", user.ID)
	}

	return user, nil
}

func (repo *UserRepository) FindByKey(ctx context.Context, key, data string) (*model.User, error) {
	var user model.User
	var query string
	
	switch key {
	case "id":
		query = `
			SELECT id, name, email, password, role, session_id, code, is_verified, created_at, updated_at, deleted_at
			FROM users
			WHERE id = $1
		`
		err := repo.Db.QueryRowContext(ctx, query, data).Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.SessionId,
			&user.Code,
			&user.IsVerified,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("user not found with id %s", data)
			}
			return nil, fmt.Errorf("error getting user by id: %w", err)
		}
	case "email":
		query = `
			SELECT id, name, email, password, role, session_id, code, is_verified, created_at, updated_at, deleted_at
			FROM users
			WHERE email = $1
		`
		err := repo.Db.QueryRowContext(ctx, query, data).Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.SessionId,
			&user.Code,
			&user.IsVerified,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("user not found with email %s", data)
			}
			return nil, fmt.Errorf("error getting user by email: %w", err)
		}
	case "username":
		query = `
			SELECT id, name, email, password, role, session_id, code, is_verified, created_at, updated_at, deleted_at
			FROM users
			WHERE name = $1
		`
		err := repo.Db.QueryRowContext(ctx, query, data).Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.SessionId,
			&user.Code,
			&user.IsVerified,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("user not found with username %s", data)
			}
			return nil, fmt.Errorf("error getting user by username: %w", err)
		}
	default:
		return nil, fmt.Errorf("unsupported key: %s", key)
	}
	
	return &user, nil
}
