package postgres

import (
	"context"
	"database/sql"
	"errors"
	"music-lib/internal/model"
	"music-lib/pkg/db"
	"time"
)

type ProfileRepository struct {
	db *db.Db
}

func NewProfileRepository(db *db.Db) *ProfileRepository {
	return &ProfileRepository{
		db: db,
	}
}

func (r *ProfileRepository) Create(ctx context.Context, entity *model.Profile) (*model.Profile, error) {
	query := `
		INSERT INTO profiles (user_id, bio, avatar_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at, updated_at
	`

	now := time.Now()
	err := r.db.QueryRowContext(ctx,
		query,
		entity.UserID,
		entity.Bio,
		entity.AvatarURL,
		now,
		now,
	).Scan(&entity.CreatedAt, &entity.UpdatedAt)

	return entity, err
}

func (r *ProfileRepository) Update(ctx context.Context, entity *model.Profile) (*model.Profile, error) {
	query := `
		UPDATE profiles 
		SET bio = $1, avatar_url = $2, updated_at = $3
		WHERE user_id = $4
		RETURNING updated_at
	`

	now := time.Now()
	err := r.db.QueryRowContext(ctx,
		query,
		entity.Bio,
		entity.AvatarURL,
		now,
		entity.UserID,
	).Scan(&entity.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}
	return entity, nil
}

func (r *ProfileRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM profiles WHERE user_id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("profile not found")
	}
	return nil
}

func (r *ProfileRepository) GetByUserID(ctx context.Context, userID uint) (*model.Profile, error) {
	var profile model.Profile
	query := `
		SELECT user_id, bio, avatar_url, created_at, updated_at
		FROM profiles
		WHERE user_id = $1
	`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&profile.UserID,
		&profile.Bio,
		&profile.AvatarURL,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}

	return &profile, nil
}
