package service

import (
	"errors"
	"fmt"
	"music-lib/internal/dto/request"
	"music-lib/internal/model"
	"music-lib/internal/repository"
	"music-lib/pkg/er"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type ProfileService struct {
	profileRepository repository.IProfileRepository
}

func NewProfileService(profile repository.IProfileRepository) *ProfileService {
	return &ProfileService{
		profileRepository: profile,
	}
}

func (s *ProfileService) NewProfile(c *gin.Context, body request.NewProfileRequest, userID uint) error {
	_, err := s.profileRepository.Create(c, &model.Profile{
		UserID:    userID,
		Bio:       body.Bio,
		AvatarURL: body.AvatarURL,
	})

	if pgErr, ok := err.(*pgconn.PgError); ok {
		if pgErr.Code == "23505" {
			return er.ErrProfileExists
		}
	}
	return err
}

func (s *ProfileService) GetProfile(c *gin.Context, userID uint) (*model.Profile, error) {
	profile, err := s.profileRepository.GetByUserID(c, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &er.NotFoundError{Message: fmt.Sprintf("profile not found for user ID %d", userID)}
		}
		return nil, &er.InternalError{Message: "failed to get profile"}
	}
	return profile, nil
}
