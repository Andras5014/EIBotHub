package service

import (
	"errors"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/model"
	"github.com/Andras5014/EIBotHub/server/internal/repository"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

type VerificationService struct {
	repo     *repository.VerificationRepository
	users    *repository.UserRepository
	activity *repository.UserActivityRepository
	hooks    *IntegrationService
}

func NewVerificationService(repo *repository.VerificationRepository, users *repository.UserRepository, activity *repository.UserActivityRepository, hooks *IntegrationService) *VerificationService {
	return &VerificationService{repo: repo, users: users, activity: activity, hooks: hooks}
}

func (s *VerificationService) Apply(userID uint, input dto.DeveloperVerificationRequest) (*dto.VerificationStatusItem, error) {
	item := &model.DeveloperVerification{
		UserID:           userID,
		VerificationType: input.VerificationType,
		RealName:         input.RealName,
		Organization:     input.Organization,
		Materials:        input.Materials,
		Reason:           input.Reason,
		Status:           "pending",
	}
	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	result := toVerificationStatus(*item)
	return &result, nil
}

func (s *VerificationService) ApplyEnterprise(userID uint, input dto.EnterpriseVerificationRequest) (*dto.VerificationStatusItem, error) {
	return s.Apply(userID, dto.DeveloperVerificationRequest{
		VerificationType: "enterprise",
		RealName:         input.RealName,
		Organization:     input.Organization,
		Materials:        input.Materials,
		Reason:           input.Reason,
	})
}

func (s *VerificationService) MyStatus(userID uint) (*dto.VerificationStatusItem, error) {
	item, err := s.repo.LatestByUser(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	result := toVerificationStatus(*item)
	return &result, nil
}

func (s *VerificationService) PublicStatus(userID uint) (*dto.VerificationStatusItem, error) {
	item, err := s.repo.LatestByUser(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	result := toVerificationStatus(*item)
	result.Materials = ""
	result.Reason = ""
	return &result, nil
}

func (s *VerificationService) AdminList() ([]dto.AdminVerificationItem, error) {
	items, err := s.repo.ListAll()
	if err != nil {
		return nil, err
	}
	result := make([]dto.AdminVerificationItem, 0, len(items))
	for _, item := range items {
		result = append(result, dto.AdminVerificationItem{
			ID:               item.ID,
			UserID:           item.UserID,
			UserName:         item.User.Username,
			VerificationType: item.VerificationType,
			RealName:         item.RealName,
			Organization:     item.Organization,
			Reason:           item.Reason,
			Status:           item.Status,
			ReviewComment:    item.ReviewComment,
			ReviewedAt:       item.ReviewedAt,
			CreatedAt:        item.CreatedAt,
		})
	}
	return result, nil
}

func (s *VerificationService) Review(id, reviewerID uint, input dto.VerificationDecisionRequest) error {
	item, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return support.NewError(http.StatusNotFound, "verification_not_found", "认证申请不存在")
		}
		return err
	}

	now := time.Now()
	item.Status = input.Decision
	item.ReviewComment = input.Comment
	item.ReviewedBy = &reviewerID
	item.ReviewedAt = &now
	if err := s.repo.Update(item); err != nil {
		return err
	}
	if input.Decision == "approved" && s.users != nil {
		if user, userErr := s.users.FindByID(item.UserID); userErr == nil {
			if user.Role == model.RoleUser {
				user.Role = model.RoleDeveloper
				_ = s.users.Update(user)
			}
		}
	}

	title := "开发者认证审核结果"
	content := "你的认证申请已审核"
	if input.Decision == "approved" {
		content = "你的认证申请已通过"
	} else {
		content = "你的认证申请已驳回"
	}
	_ = s.activity.AddNotification(&model.Notification{
		UserID:  item.UserID,
		Type:    "verification_review",
		Title:   title,
		Content: content,
	})
	if s.hooks != nil {
		s.hooks.Emit(item.UserID, WebhookEventVerificationReviewed, map[string]any{
			"verification_id":    item.ID,
			"verification_type":  item.VerificationType,
			"status":             item.Status,
			"review_comment":     item.ReviewComment,
			"reviewed_at":        item.ReviewedAt,
			"reviewer_user_id":   reviewerID,
			"applicant_user_id":  item.UserID,
			"applicant_username": item.User.Username,
		})
	}
	return nil
}

func toVerificationStatus(item model.DeveloperVerification) dto.VerificationStatusItem {
	return dto.VerificationStatusItem{
		ID:               item.ID,
		UserID:           item.UserID,
		VerificationType: item.VerificationType,
		RealName:         item.RealName,
		Organization:     item.Organization,
		Materials:        item.Materials,
		Reason:           item.Reason,
		Status:           item.Status,
		ReviewComment:    item.ReviewComment,
		ReviewedAt:       item.ReviewedAt,
		CreatedAt:        item.CreatedAt,
	}
}
