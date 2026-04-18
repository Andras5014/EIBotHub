package repository

import (
	"gorm.io/gorm"

	"github.com/Andras5014/EIBotHub/server/internal/model"
)

type RewardRepository struct {
	db *gorm.DB
}

func NewRewardRepository(db *gorm.DB) *RewardRepository {
	return &RewardRepository{db: db}
}

func (r *RewardRepository) Add(item *model.RewardLedger) error {
	return r.db.Create(item).Error
}

func (r *RewardRepository) InTx(fn func(repo *RewardRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(NewRewardRepository(tx))
	})
}

func (r *RewardRepository) SumPoints(userID uint) (int64, error) {
	var total int64
	err := r.db.Model(&model.RewardLedger{}).Where("user_id = ?", userID).Select("coalesce(sum(points),0)").Scan(&total).Error
	return total, err
}

func (r *RewardRepository) ListByUser(userID uint) ([]model.RewardLedger, error) {
	var items []model.RewardLedger
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&items).Error
	return items, err
}

func (r *RewardRepository) Rankings(limit int) ([]struct {
	UserID   uint
	UserName string
	Points   int64
}, error) {
	var rows []struct {
		UserID   uint
		UserName string
		Points   int64
	}
	err := r.db.Model(&model.RewardLedger{}).
		Select("reward_ledgers.user_id as user_id, users.username as user_name, sum(reward_ledgers.points) as points").
		Joins("left join users on users.id = reward_ledgers.user_id").
		Group("reward_ledgers.user_id, users.username").
		Order("points desc").
		Limit(limit).
		Scan(&rows).Error
	return rows, err
}

func (r *RewardRepository) ListBenefits() ([]model.RewardBenefit, error) {
	var items []model.RewardBenefit
	err := r.db.Where("active = ?", true).Order("cost_points asc").Find(&items).Error
	return items, err
}

func (r *RewardRepository) ListAllBenefits() ([]model.RewardBenefit, error) {
	var items []model.RewardBenefit
	err := r.db.Order("cost_points asc, id asc").Find(&items).Error
	return items, err
}

func (r *RewardRepository) FindBenefit(id uint) (*model.RewardBenefit, error) {
	var item model.RewardBenefit
	err := r.db.First(&item, id).Error
	return &item, err
}

func (r *RewardRepository) UpdateBenefit(item *model.RewardBenefit) error {
	return r.db.Save(item).Error
}

func (r *RewardRepository) AddRedemption(item *model.RewardRedemption) error {
	return r.db.Create(item).Error
}

func (r *RewardRepository) ListRedemptions(userID uint) ([]model.RewardRedemption, error) {
	var items []model.RewardRedemption
	err := r.db.Preload("Benefit").Where("user_id = ?", userID).Order("created_at desc").Find(&items).Error
	return items, err
}

func (r *RewardRepository) CountBenefits(activeOnly bool) (int64, error) {
	var count int64
	db := r.db.Model(&model.RewardBenefit{})
	if activeOnly {
		db = db.Where("active = ?", true)
	}
	err := db.Count(&count).Error
	return count, err
}

func (r *RewardRepository) CountRedemptions() (int64, error) {
	var count int64
	err := r.db.Model(&model.RewardRedemption{}).Count(&count).Error
	return count, err
}

func (r *RewardRepository) CountLedgerEntries() (int64, error) {
	var count int64
	err := r.db.Model(&model.RewardLedger{}).Count(&count).Error
	return count, err
}

func (r *RewardRepository) SumAllPoints() (int64, error) {
	var total int64
	err := r.db.Model(&model.RewardLedger{}).Select("coalesce(sum(points),0)").Scan(&total).Error
	return total, err
}

func (r *RewardRepository) ListAdminAdjustments(limit int) ([]model.RewardLedger, error) {
	var items []model.RewardLedger
	err := r.db.Preload("User").
		Where("source_type = ?", "admin_adjustment").
		Order("created_at desc").
		Limit(limit).
		Find(&items).Error
	return items, err
}
