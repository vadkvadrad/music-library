package service

import (
	"music-lib/internal/model"
	"music-lib/internal/repository"
	"music-lib/pkg/er"

	"github.com/gin-gonic/gin"
)

type PermissionService struct {
	permissionRepo repository.IPermissionRepository
}

func NewPermissionService(permission repository.IPermissionRepository) *PermissionService {
	return &PermissionService{
		permissionRepo: permission,
	}
}

func (s *PermissionService) HasPermission(
	userID, resourceID uint,
	resourceType model.Resource,
	permission model.Permission,
) bool {
	return s.permissionRepo.HasPermission(userID, resourceID, resourceType, permission)
}

func (s *PermissionService) AddPermission(
	ctx *gin.Context,
	userID, resourceID uint,
	resourceType model.Resource,
	permission model.Permission,
) error {
	_, err := s.permissionRepo.Create(ctx, &model.ResourcePermission{
		UserID: userID,
		ResourceID: resourceID,
		ResourceType: resourceType,
		Permission: permission,
	})
	if err != nil {
		return er.InternalError{Message: err.Error()}
	}
	return nil
}