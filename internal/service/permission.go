package service

import (
	"music-lib/internal/model"
	"music-lib/internal/repository"
	"music-lib/pkg/er"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type PermissionService struct {
	permissionRepo repository.IPermissionRepository

	logger *zap.SugaredLogger
}

func NewPermissionService(permission repository.IPermissionRepository, log *zap.SugaredLogger) *PermissionService {
	return &PermissionService{
		permissionRepo: permission,
		logger: log,
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
	s.logger.Debugw("Adding new permission",
		"user id", userID,
		"resource id", resourceID,
		"resource type", resourceType,
		"permission", permission,
	)
	_, err := s.permissionRepo.Create(ctx, &model.ResourcePermission{
		UserID: userID,
		ResourceID: resourceID,
		ResourceType: resourceType,
		Permission: permission,
	})
	if err != nil {
		s.logger.Errorw("failed to add permission",
			"error", err.Error(),
			"error type", "Internal",
			"user id", userID,
			"resource id", resourceID,
			"resource type", resourceType,
			"permission", permission,
		)	
		return er.InternalError{Message: err.Error()}
	}
	s.logger.Debugw("Permission added successfully")
	return nil
}