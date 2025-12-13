package postgres

import (
	"context"
	"database/sql"
	"errors"
	"music-lib/internal/model"
	"music-lib/pkg/db"
	"time"
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
	query := `
		INSERT INTO resource_permissions (user_id, resource_id, resource_type, permission, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	now := time.Now()
	err := r.db.QueryRowContext(ctx,
		query,
		entity.UserID,
		entity.ResourceID,
		entity.ResourceType,
		entity.Permission,
		now,
		now,
	).Scan(&entity.ID, &entity.CreatedAt, &entity.UpdatedAt)

	return entity, err
}

func (r *PermissionRepository) Update(ctx context.Context, entity *model.ResourcePermission) (*model.ResourcePermission, error) {
	query := `
		UPDATE resource_permissions 
		SET user_id = $1, resource_id = $2, resource_type = $3, permission = $4, updated_at = $5
		WHERE id = $6
		RETURNING updated_at
	`

	now := time.Now()
	err := r.db.QueryRowContext(ctx,
		query,
		entity.UserID,
		entity.ResourceID,
		entity.ResourceType,
		entity.Permission,
		now,
		entity.ID,
	).Scan(&entity.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("permission not found")
		}
		return nil, err
	}
	return entity, nil
}

func (r *PermissionRepository) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM resource_permissions WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("permission not found")
	}
	return nil
}

func (r *PermissionRepository) HasPermission(
	userID, resourceID uint,
	resourceType model.Resource,
	permission model.Permission,
) bool {
	var count int64
	query := `
		SELECT COUNT(*) 
		FROM resource_permissions
		WHERE user_id = $1 AND resource_id = $2 AND resource_type = $3 AND permission = $4
	`

	err := r.db.QueryRow(query, userID, resourceID, resourceType, permission).Scan(&count)
	return err == nil && count > 0
}
