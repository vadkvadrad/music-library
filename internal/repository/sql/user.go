package sql

import (
	"context"
	"database/sql"
	"fmt"
	"music-lib/internal/model"
)

type UserRepository struct {
	db *DB
}

func NewUserRepository(db *DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) Create(ctx context.Context, entity *model.User) (*model.User, error) {
	query := `
		INSERT INTO users (name, email, password, role, session_id, code, is_verified, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	var id uint
	err := r.db.QueryRowContext(ctx, query,
		entity.Name,
		entity.Email,
		entity.Password,
		entity.Role,
		entity.SessionId,
		entity.Code,
		entity.IsVerified,
		entity.CreatedAt,
		entity.UpdatedAt,
		entity.DeletedAt,
	).Scan(&id)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	entity.ID = id
	return entity, nil
}

func (r *UserRepository) Update(ctx context.Context, entity *model.User) (*model.User, error) {
	query := `
		UPDATE users
		SET name = $1, email = $2, password = $3, role = $4, session_id = $5, code = $6, is_verified = $7, updated_at = $8
		WHERE id = $9
	`

	result, err := r.db.ExecContext(ctx, query,
		entity.Name,
		entity.Email,
		entity.Password,
		entity.Role,
		entity.SessionId,
		entity.Code,
		entity.IsVerified,
		entity.UpdatedAt,
		entity.ID,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return nil, fmt.Errorf("no user found with id %d", entity.ID)
	}

	return entity, nil
}

func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no user found with id %d", id)
	}

	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	query := `
		SELECT id, name, email, password, role, session_id, code, is_verified, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1
	`
	
	var user model.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
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
			return nil, fmt.Errorf("user not found with id %d", id)
		}
		return nil, fmt.Errorf("error getting user by id: %w", err)
	}
	
	return &user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, name, email, password, role, session_id, code, is_verified, created_at, updated_at, deleted_at
		FROM users
		WHERE email = $1
	`
	
	var user model.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
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
			return nil, fmt.Errorf("user not found with email %s", email)
		}
		return nil, fmt.Errorf("error getting user by email: %w", err)
	}
	
	return &user, nil
}

func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT id, name, email, password, role, session_id, code, is_verified, created_at, updated_at, deleted_at
		FROM users
		WHERE name = $1
	`
	
	var user model.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
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
			return nil, fmt.Errorf("user not found with username %s", username)
		}
		return nil, fmt.Errorf("error getting user by username: %w", err)
	}
	
	return &user, nil
}