package service

import (
	"net/http"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/model"
	"github.com/Andras5014/EIBotHub/server/internal/repository"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

type UserService struct {
	users     *repository.UserRepository
	activity  *repository.UserActivityRepository
	models    *ModelService
	datasets  *DatasetService
}

func NewUserService(users *repository.UserRepository, activity *repository.UserActivityRepository, models *ModelService, datasets *DatasetService) *UserService {
	return &UserService{users: users, activity: activity, models: models, datasets: datasets}
}

func (s *UserService) Profile(userID uint) (*dto.UserSummary, error) {
	user, err := s.users.FindByID(userID)
	if err != nil {
		return nil, support.NewError(http.StatusNotFound, "user_not_found", "user not found")
	}
	summary := toUserSummary(*user)
	return &summary, nil
}

func (s *UserService) PublicProfile(userID uint) (*dto.UserSummary, error) {
	user, err := s.users.FindByID(userID)
	if err != nil {
		return nil, support.NewError(http.StatusNotFound, "user_not_found", "user not found")
	}
	summary := toUserSummary(*user)
	summary.Email = ""
	return &summary, nil
}

func (s *UserService) UpdateProfile(userID uint, input dto.ProfileUpdateRequest) (*dto.UserSummary, error) {
	user, err := s.users.FindByID(userID)
	if err != nil {
		return nil, support.NewError(http.StatusNotFound, "user_not_found", "user not found")
	}
	user.Username = input.Username
	user.Bio = input.Bio
	user.Avatar = input.Avatar
	if err := s.users.Update(user); err != nil {
		return nil, err
	}
	summary := toUserSummary(*user)
	return &summary, nil
}

func (s *UserService) Favorites(userID uint) ([]model.Favorite, error) {
	return s.activity.ListFavorites(userID)
}

func (s *UserService) ToggleFavorite(userID uint, input dto.FavoriteRequest) error {
	favorited, err := s.activity.IsFavorited(userID, input.ResourceType, input.ResourceID)
	if err != nil {
		return err
	}
	if favorited {
		return s.activity.RemoveFavorite(userID, input.ResourceType, input.ResourceID)
	}
	return s.activity.AddFavorite(&model.Favorite{
		UserID:        userID,
		ResourceType:  input.ResourceType,
		ResourceID:    input.ResourceID,
		ResourceTitle: input.Title,
	})
}

func (s *UserService) Downloads(userID uint) ([]model.DownloadRecord, error) {
	return s.activity.ListDownloads(userID)
}

func (s *UserService) Notifications(userID uint) ([]model.Notification, error) {
	return s.activity.ListNotifications(userID)
}

func (s *UserService) MarkNotificationsRead(userID uint) error {
	return s.activity.MarkNotificationsRead(userID)
}

func (s *UserService) Uploads(userID uint) (map[string]any, error) {
	models, err := s.models.Mine(userID)
	if err != nil {
		return nil, err
	}
	datasets, err := s.datasets.Mine(userID)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"models":   models,
		"datasets": datasets,
	}, nil
}
