package postgres

import (
	"context"
	"music-lib/internal/model"
	"music-lib/pkg/db"
)

type ProfileRepository struct {
	db *db.Db
}

func NewProfileRepository(db *db.Db) *ProfileRepository {
	return &ProfileRepository{
		db: db,
	}
}


func (r *ProfileRepository) Create(ctx context.Context, entity *model.Profile) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *ProfileRepository) Update(ctx context.Context, entity *model.Profile) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *ProfileRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Profile{}, id).Error
}

func (r *ProfileRepository) GetByUserID(ctx context.Context, userID uint) (*model.Profile, error) {
	var profile model.Profile
    err := r.db.WithContext(ctx).
		Preload("Collections").
        Where("user_id = ?", userID).
        First(&profile).
        Error

    if err != nil {
		return nil, err
    }

    return &profile, nil
}