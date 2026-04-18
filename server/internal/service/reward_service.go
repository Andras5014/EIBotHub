package service

import (
	"net/http"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/model"
	"github.com/Andras5014/EIBotHub/server/internal/repository"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

type RewardService struct {
	repo     *repository.RewardRepository
	users    *repository.UserRepository
	activity *repository.UserActivityRepository
	hooks    *IntegrationService
}

func NewRewardService(repo *repository.RewardRepository, users *repository.UserRepository, activity *repository.UserActivityRepository, hooks *IntegrationService) *RewardService {
	return &RewardService{repo: repo, users: users, activity: activity, hooks: hooks}
}

func (s *RewardService) Add(userID uint, sourceType string, points int, remark string) {
	if s == nil || s.repo == nil || userID == 0 || points == 0 {
		return
	}
	_ = s.repo.Add(&model.RewardLedger{
		UserID:     userID,
		SourceType: sourceType,
		Points:     points,
		Remark:     remark,
	})
}

func (s *RewardService) Summary(userID uint) (*dto.RewardSummary, error) {
	points, err := s.repo.SumPoints(userID)
	if err != nil {
		return nil, err
	}
	return &dto.RewardSummary{Points: points}, nil
}

func (s *RewardService) Ledger(userID uint) ([]dto.RewardLedgerItem, error) {
	items, err := s.repo.ListByUser(userID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.RewardLedgerItem, 0, len(items))
	for _, item := range items {
		result = append(result, dto.RewardLedgerItem{
			ID:         item.ID,
			SourceType: item.SourceType,
			Points:     item.Points,
			Remark:     item.Remark,
			CreatedAt:  item.CreatedAt,
		})
	}
	return result, nil
}

func (s *RewardService) Rankings() ([]dto.ContributorRankingItem, error) {
	rows, err := s.repo.Rankings(10)
	if err != nil {
		return nil, err
	}
	result := make([]dto.ContributorRankingItem, 0, len(rows))
	for _, row := range rows {
		result = append(result, dto.ContributorRankingItem{
			UserID:   row.UserID,
			UserName: row.UserName,
			Points:   row.Points,
		})
	}
	return result, nil
}

func (s *RewardService) Benefits() ([]dto.RewardBenefitItem, error) {
	items, err := s.repo.ListBenefits()
	if err != nil {
		return nil, err
	}
	result := make([]dto.RewardBenefitItem, 0, len(items))
	for _, item := range items {
		result = append(result, dto.RewardBenefitItem{
			ID:         item.ID,
			Name:       item.Name,
			Summary:    item.Summary,
			CostPoints: item.CostPoints,
			Active:     item.Active,
		})
	}
	return result, nil
}

func (s *RewardService) Redeem(userID uint, benefitID uint) (*dto.RewardRedemptionItem, error) {
	benefit, err := s.repo.FindBenefit(benefitID)
	if err != nil {
		return nil, support.NewError(http.StatusNotFound, "benefit_not_found", "兑换权益不存在")
	}
	if !benefit.Active {
		return nil, support.NewError(http.StatusBadRequest, "benefit_inactive", "兑换权益不可用")
	}
	points, err := s.repo.SumPoints(userID)
	if err != nil {
		return nil, err
	}
	if points < int64(benefit.CostPoints) {
		return nil, support.NewError(http.StatusBadRequest, "points_insufficient", "积分不足，无法兑换")
	}
	redemption := &model.RewardRedemption{
		UserID:     userID,
		BenefitID:  benefit.ID,
		CostPoints: benefit.CostPoints,
	}
	if err := s.repo.InTx(func(repo *repository.RewardRepository) error {
		if err := repo.AddRedemption(redemption); err != nil {
			return err
		}
		return repo.Add(&model.RewardLedger{
			UserID:     userID,
			SourceType: "redeem",
			Points:     -benefit.CostPoints,
			Remark:     "兑换权益：" + benefit.Name,
		})
	}); err != nil {
		return nil, err
	}
	return &dto.RewardRedemptionItem{
		ID:          redemption.ID,
		BenefitID:   benefit.ID,
		BenefitName: benefit.Name,
		CostPoints:  benefit.CostPoints,
		CreatedAt:   redemption.CreatedAt,
	}, nil
}

func (s *RewardService) Redemptions(userID uint) ([]dto.RewardRedemptionItem, error) {
	items, err := s.repo.ListRedemptions(userID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.RewardRedemptionItem, 0, len(items))
	for _, item := range items {
		result = append(result, dto.RewardRedemptionItem{
			ID:          item.ID,
			BenefitID:   item.BenefitID,
			BenefitName: item.Benefit.Name,
			CostPoints:  item.CostPoints,
			CreatedAt:   item.CreatedAt,
		})
	}
	return result, nil
}

func (s *RewardService) EmitRedeemed(userID uint, item dto.RewardRedemptionItem) {
	if s == nil || s.hooks == nil {
		return
	}
	s.hooks.Emit(userID, WebhookEventRewardRedeemed, map[string]any{
		"redemption_id": item.ID,
		"benefit_id":    item.BenefitID,
		"benefit_name":  item.BenefitName,
		"cost_points":   item.CostPoints,
		"created_at":    item.CreatedAt,
	})
}

func (s *RewardService) AdminOverview() (*dto.AdminRewardOverview, error) {
	benefits, err := s.repo.CountBenefits(false)
	if err != nil {
		return nil, err
	}
	activeBenefits, err := s.repo.CountBenefits(true)
	if err != nil {
		return nil, err
	}
	redemptions, err := s.repo.CountRedemptions()
	if err != nil {
		return nil, err
	}
	ledgerEntries, err := s.repo.CountLedgerEntries()
	if err != nil {
		return nil, err
	}
	netPoints, err := s.repo.SumAllPoints()
	if err != nil {
		return nil, err
	}
	return &dto.AdminRewardOverview{
		Benefits:       benefits,
		ActiveBenefits: activeBenefits,
		Redemptions:    redemptions,
		LedgerEntries:  ledgerEntries,
		NetPoints:      netPoints,
	}, nil
}

func (s *RewardService) AdminBenefits() ([]dto.RewardBenefitItem, error) {
	items, err := s.repo.ListAllBenefits()
	if err != nil {
		return nil, err
	}
	result := make([]dto.RewardBenefitItem, 0, len(items))
	for _, item := range items {
		result = append(result, dto.RewardBenefitItem{
			ID:         item.ID,
			Name:       item.Name,
			Summary:    item.Summary,
			CostPoints: item.CostPoints,
			Active:     item.Active,
		})
	}
	return result, nil
}

func (s *RewardService) UpdateBenefit(id uint, input dto.AdminRewardBenefitRequest) (*dto.RewardBenefitItem, error) {
	item, err := s.repo.FindBenefit(id)
	if err != nil {
		return nil, support.NewError(http.StatusNotFound, "benefit_not_found", "兑换权益不存在")
	}
	item.Name = input.Name
	item.Summary = input.Summary
	item.CostPoints = input.CostPoints
	item.Active = input.Active
	if err := s.repo.UpdateBenefit(item); err != nil {
		return nil, err
	}
	return &dto.RewardBenefitItem{
		ID:         item.ID,
		Name:       item.Name,
		Summary:    item.Summary,
		CostPoints: item.CostPoints,
		Active:     item.Active,
	}, nil
}

func (s *RewardService) AdjustPoints(input dto.AdminRewardAdjustmentRequest) (*dto.AdminRewardAdjustmentItem, error) {
	user, err := s.users.FindByID(input.UserID)
	if err != nil {
		return nil, support.NewError(http.StatusNotFound, "user_not_found", "用户不存在")
	}
	entry := &model.RewardLedger{
		UserID:     input.UserID,
		SourceType: "admin_adjustment",
		Points:     input.Points,
		Remark:     input.Remark,
	}
	if err := s.repo.Add(entry); err != nil {
		return nil, err
	}
	if s.activity != nil {
		direction := "增加"
		if input.Points < 0 {
			direction = "扣减"
		}
		_ = s.activity.AddNotification(&model.Notification{
			UserID:  input.UserID,
			Type:    "reward_adjustment",
			Title:   "积分变更通知",
			Content: "管理员已" + direction + "你的积分：" + input.Remark,
		})
	}
	return &dto.AdminRewardAdjustmentItem{
		ID:        entry.ID,
		UserID:    entry.UserID,
		UserName:  user.Username,
		Points:    entry.Points,
		Remark:    entry.Remark,
		CreatedAt: entry.CreatedAt,
	}, nil
}

func (s *RewardService) AdminAdjustments() ([]dto.AdminRewardAdjustmentItem, error) {
	items, err := s.repo.ListAdminAdjustments(20)
	if err != nil {
		return nil, err
	}
	result := make([]dto.AdminRewardAdjustmentItem, 0, len(items))
	for _, item := range items {
		result = append(result, dto.AdminRewardAdjustmentItem{
			ID:        item.ID,
			UserID:    item.UserID,
			UserName:  item.User.Username,
			Points:    item.Points,
			Remark:    item.Remark,
			CreatedAt: item.CreatedAt,
		})
	}
	return result, nil
}
