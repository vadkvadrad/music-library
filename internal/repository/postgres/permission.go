package postgres

import (
	"context"
	"music-lib/internal/model"
	"music-lib/pkg/db"
)

type PermissionRepository struct {
	db *db.Db
}

func NewPermissionRepository(db *db.Db) *PermissionRepository {
	return &PermissionRepository{
		db: db,
	}
}

func (r *PermissionRepository) Create(ctx context.Context, entity *model.ResourcePermission) (*model.ResourcePermission, error) {
	err := r.db.WithContext(ctx).Create(entity).Error
	return entity, err
}

func (r *PermissionRepository) Update(ctx context.Context, entity *model.ResourcePermission) (*model.ResourcePermission, error) {
	err := r.db.WithContext(ctx).Save(entity).Error
	return entity, err
}

func (r *PermissionRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.ResourcePermission{}, id).Error
}

func (r *PermissionRepository) HasPermission(
	userID, resourceID uint,
	resourceType model.Resource,
	permission model.Permission,
) bool {
	var count int64
	err := r.db.Find(&model.ResourcePermission{}).
		Where("user_id = ? AND resource_id = ? AND resource_type = ? AND permission = ?",
			userID, resourceID, resourceType, permission).
		Count(&count).Error
	return err == nil && count > 0
}
