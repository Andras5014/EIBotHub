package repository

import (
	"gorm.io/gorm"

	"github.com/Andras5014/EIBotHub/server/internal/model"
)

type WikiRepository struct {
	db *gorm.DB
}

func NewWikiRepository(db *gorm.DB) *WikiRepository {
	return &WikiRepository{db: db}
}

func (r *WikiRepository) CreatePage(item *model.WikiPage) error {
	return r.db.Create(item).Error
}

func (r *WikiRepository) UpdatePage(item *model.WikiPage) error {
	return r.db.Save(item).Error
}

func (r *WikiRepository) AddRevision(item *model.WikiRevision) error {
	return r.db.Create(item).Error
}

func (r *WikiRepository) ListPages() ([]model.WikiPage, error) {
	var items []model.WikiPage
	err := r.db.Preload("Editor").Where("status = ?", model.StatusPublished).Order("updated_at desc").Find(&items).Error
	return items, err
}

func (r *WikiRepository) ListAllPages() ([]model.WikiPage, error) {
	var items []model.WikiPage
	err := r.db.Preload("Editor").Order("updated_at desc").Find(&items).Error
	return items, err
}

func (r *WikiRepository) GetPage(id uint) (*model.WikiPage, error) {
	var item model.WikiPage
	err := r.db.Preload("Editor").First(&item, id).Error
	return &item, err
}

func (r *WikiRepository) ListRevisions(pageID uint) ([]model.WikiRevision, error) {
	var items []model.WikiRevision
	err := r.db.Where("page_id = ?", pageID).Order("created_at desc").Find(&items).Error
	return items, err
}

func (r *WikiRepository) GetRevision(pageID, revisionID uint) (*model.WikiRevision, error) {
	var item model.WikiRevision
	err := r.db.Where("page_id = ? AND id = ?", pageID, revisionID).First(&item).Error
	return &item, err
}

func (r *WikiRepository) CountRevisions(pageID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.WikiRevision{}).Where("page_id = ?", pageID).Count(&count).Error
	return count, err
}
