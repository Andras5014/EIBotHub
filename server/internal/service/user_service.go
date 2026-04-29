package service

import (
	"errors"
	"net/http"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/model"
	"github.com/Andras5014/EIBotHub/server/internal/repository"
	"github.com/Andras5014/EIBotHub/server/internal/support"
	"gorm.io/gorm"
)

type UserService struct {
	users    *repository.UserRepository
	activity *repository.UserActivityRepository
	models   *ModelService
	datasets *DatasetService
	reviews  *repository.ReviewRepository
}

func NewUserService(users *repository.UserRepository, activity *repository.UserActivityRepository, models *ModelService, datasets *DatasetService, reviews *repository.ReviewRepository) *UserService {
	return &UserService{users: users, activity: activity, models: models, datasets: datasets, reviews: reviews}
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
	summary.Permissions = nil
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

func (s *UserService) MarkNotificationRead(userID, notificationID uint) error {
	if err := s.activity.MarkNotificationRead(userID, notificationID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return support.NewError(http.StatusNotFound, "notification_not_found", "notification not found")
		}
		return err
	}
	return nil
}

func (s *UserService) MarkNotificationsRead(userID uint) error {
	return s.activity.MarkNotificationsRead(userID)
}

func (s *UserService) Uploads(userID uint) (map[string]any, error) {
	models, err := s.models.Mine(userID)
	if err != nil {
		return nil, err
	}
	if err := s.attachReviewComments(models, "models"); err != nil {
		return nil, err
	}
	datasets, err := s.datasets.Mine(userID)
	if err != nil {
		return nil, err
	}
	if err := s.attachReviewComments(datasets, "datasets"); err != nil {
		return nil, err
	}
	return map[string]any{
		"models":   models,
		"datasets": datasets,
	}, nil
}

func (s *UserService) attachReviewComments(items []dto.ResourceCard, reviewType string) error {
	if s.reviews == nil || len(items) == 0 {
		return nil
	}
	ids := make([]uint, 0, len(items))
	for _, item := range items {
		ids = append(ids, item.ID)
	}
	latestLogs, err := s.reviews.ListLatestByResourceIDs(reviewType, ids)
	if err != nil {
		return err
	}
	for index := range items {
		if items[index].Status != model.StatusRejected {
			continue
		}
		if logItem, ok := latestLogs[items[index].ID]; ok && logItem.Decision == "rejected" {
			items[index].ReviewComment = logItem.Comment
		}
	}
	return nil
}

func latestReviewComment(reviews *repository.ReviewRepository, reviewType string, resourceID uint, status string, isOwner bool) (string, bool, error) {
	if reviews == nil || !isOwner || status != model.StatusRejected {
		return "", false, nil
	}
	latestLogs, err := reviews.ListLatestByResourceIDs(reviewType, []uint{resourceID})
	if err != nil {
		return "", false, err
	}
	logItem, ok := latestLogs[resourceID]
	if !ok || logItem.Decision != "rejected" {
		return "", false, nil
	}
	return logItem.Comment, true, nil
}
