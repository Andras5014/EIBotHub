package repository

import (
	"gorm.io/gorm"

	"github.com/Andras5014/EIBotHub/server/internal/model"
)

type VerificationRepository struct {
	db *gorm.DB
}

func NewVerificationRepository(db *gorm.DB) *VerificationRepository {
	return &VerificationRepository{db: db}
}

func (r *VerificationRepository) Create(item *model.DeveloperVerification) error {
	return r.db.Create(item).Error
}

func (r *VerificationRepository) Update(item *model.DeveloperVerification) error {
	return r.db.Save(item).Error
}

func (r *VerificationRepository) LatestByUser(userID uint) (*model.DeveloperVerification, error) {
	var item model.DeveloperVerification
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").First(&item).Error
	return &item, err
}

func (r *VerificationRepository) FindByID(id uint) (*model.DeveloperVerification, error) {
	var item model.DeveloperVerification
	err := r.db.Preload("User").First(&item, id).Error
	return &item, err
}

func (r *VerificationRepository) ListAll() ([]model.DeveloperVerification, error) {
	var items []model.DeveloperVerification
	err := r.db.Preload("User").Order("created_at desc").Find(&items).Error
	return items, err
}
