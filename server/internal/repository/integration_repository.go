package repository

import (
	"gorm.io/gorm"

	"github.com/Andras5014/EIBotHub/server/internal/model"
)

type IntegrationRepository struct {
	db *gorm.DB
}

func NewIntegrationRepository(db *gorm.DB) *IntegrationRepository {
	return &IntegrationRepository{db: db}
}

func (r *IntegrationRepository) CreateWebhook(item *model.WebhookSubscription) error {
	return r.db.Create(item).Error
}

func (r *IntegrationRepository) UpdateWebhook(item *model.WebhookSubscription) error {
	return r.db.Save(item).Error
}

func (r *IntegrationRepository) ListWebhooksByUser(userID uint) ([]model.WebhookSubscription, error) {
	var items []model.WebhookSubscription
	err := r.db.Where("user_id = ?", userID).Order("updated_at desc, created_at desc").Find(&items).Error
	return items, err
}

func (r *IntegrationRepository) GetWebhookByIDAndUser(id, userID uint) (*model.WebhookSubscription, error) {
	var item model.WebhookSubscription
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&item).Error
	return &item, err
}

func (r *IntegrationRepository) AddDelivery(item *model.WebhookDelivery) error {
	return r.db.Create(item).Error
}

func (r *IntegrationRepository) ListDeliveries(webhookID uint) ([]model.WebhookDelivery, error) {
	var items []model.WebhookDelivery
	err := r.db.Where("webhook_id = ?", webhookID).Order("created_at desc").Limit(20).Find(&items).Error
	return items, err
}
