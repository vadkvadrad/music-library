package service

import (
	"fmt"
	"music-lib/internal/dto/request"
	"music-lib/internal/model"
	"music-lib/internal/repository"
	"music-lib/pkg/er"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
)

type ProfileService struct {
	profileRepository repository.IProfileRepository
	artistRepository  repository.IArtistRepository
	userRepository    repository.IUserRepository
}

func NewProfileService(profile repository.IProfileRepository, artist repository.IArtistRepository, user repository.IUserRepository) *ProfileService {
	return &ProfileService{
		profileRepository: profile,
		artistRepository:  artist,
		userRepository:    user,
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
		if err.Error() == "profile not found" {
			return nil, &er.NotFoundError{Message: fmt.Sprintf("profile not found for user ID %d", userID)}
		}
		return nil, &er.InternalError{Message: "failed to get profile"}
	}
	return profile, nil
}

func (s *ProfileService) GetArtistByUserID(c *gin.Context, userID uint) (*model.Artist, error) {
	artist, err := s.artistRepository.GetByUserID(c, userID)
	if err != nil {
		// Если артист не найден, это не ошибка - просто возвращаем nil
		return nil, nil
	}
	return artist, nil
}

func (s *ProfileService) GetUserByID(userID uint) (*model.User, error) {
	user, err := s.userRepository.FindByKey("id", fmt.Sprintf("%d", userID))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *ProfileService) UpdateProfile(c *gin.Context, userID uint, body request.UpdateProfileRequest) (*model.Profile, error) {
	profile, err := s.profileRepository.GetByUserID(c, userID)
	if err != nil {
		if err.Error() == "profile not found" {
			return nil, &er.NotFoundError{Message: fmt.Sprintf("profile not found for user ID %d", userID)}
		}
		return nil, &er.InternalError{Message: "failed to get profile"}
	}

	if body.Bio != "" {
		profile.Bio = body.Bio
	}

	if body.AvatarURL != "" {
		profile.AvatarURL = body.AvatarURL
	}

	return s.profileRepository.Update(c, profile)
}
