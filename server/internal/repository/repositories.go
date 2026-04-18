package repository

import (
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/Andras5014/EIBotHub/server/internal/model"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", strings.ToLower(email)).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Count() (int64, error) {
	var count int64
	return count, r.db.Model(&model.User{}).Count(&count).Error
}

func (r *UserRepository) SearchByKeyword(query string) ([]model.User, error) {
	var users []model.User
	like := "%" + strings.TrimSpace(query) + "%"
	err := r.db.Where("username LIKE ? OR email LIKE ? OR bio LIKE ?", like, like, like).
		Order("updated_at desc").
		Find(&users).Error
	return users, err
}

type PortalRepository struct {
	db *gorm.DB
}

func NewPortalRepository(db *gorm.DB) *PortalRepository {
	return &PortalRepository{db: db}
}

func (r *PortalRepository) ListAnnouncements() ([]model.Announcement, error) {
	var announcements []model.Announcement
	err := r.db.Order("pinned desc, published_at desc").Find(&announcements).Error
	return announcements, err
}

func (r *PortalRepository) CreateAnnouncement(item *model.Announcement) error {
	return r.db.Create(item).Error
}

func (r *PortalRepository) CountAnnouncements() (int64, error) {
	var count int64
	return count, r.db.Model(&model.Announcement{}).Count(&count).Error
}

func (r *PortalRepository) ListModuleSettings() ([]model.HomeModuleSetting, error) {
	var items []model.HomeModuleSetting
	err := r.db.Order("module_key asc").Find(&items).Error
	return items, err
}

func (r *PortalRepository) UpsertModuleSetting(item *model.HomeModuleSetting) error {
	var existing model.HomeModuleSetting
	err := r.db.Where("module_key = ?", item.ModuleKey).First(&existing).Error
	if err == nil {
		existing.Enabled = item.Enabled
		return r.db.Save(&existing).Error
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	return r.db.Create(item).Error
}

func (r *PortalRepository) ListFeaturedResources() ([]model.FeaturedResource, error) {
	var items []model.FeaturedResource
	err := r.db.Order("sort_order asc, updated_at desc").Find(&items).Error
	return items, err
}

func (r *PortalRepository) CreateFeaturedResource(item *model.FeaturedResource) error {
	return r.db.Create(item).Error
}

func (r *PortalRepository) GetFeaturedResource(id uint) (*model.FeaturedResource, error) {
	var item model.FeaturedResource
	err := r.db.First(&item, id).Error
	return &item, err
}

func (r *PortalRepository) UpdateFeaturedResource(item *model.FeaturedResource) error {
	return r.db.Save(item).Error
}

func (r *PortalRepository) DeleteFeaturedResource(id uint) error {
	return r.db.Delete(&model.FeaturedResource{}, id).Error
}

type ModelRepository struct {
	db *gorm.DB
}

func NewModelRepository(db *gorm.DB) *ModelRepository {
	return &ModelRepository{db: db}
}

func (r *ModelRepository) Create(item *model.ModelAsset) error {
	return r.db.Create(item).Error
}

func (r *ModelRepository) Update(item *model.ModelAsset) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(item).Error
}

func (r *ModelRepository) GetByID(id uint) (*model.ModelAsset, error) {
	var item model.ModelAsset
	err := r.db.Preload("Owner").Preload("Versions", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at desc")
	}).First(&item, id).Error
	return &item, err
}

func (r *ModelRepository) ListPublished(query, tags, robotType, sort string, page, pageSize int) ([]model.ModelAsset, int64, error) {
	db := r.db.Model(&model.ModelAsset{}).Where("status = ?", model.StatusPublished)
	db = applyTextSearch(db, query, []string{"name", "summary", "description"})
	db = applyTagFilter(db, tags)
	if robotType != "" {
		db = db.Where("robot_type = ?", robotType)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []model.ModelAsset
	err := applyPagination(applySort(db, sort, "downloads desc, updated_at desc"), page, pageSize).
		Preload("Owner").
		Find(&items).Error
	return items, total, err
}

func (r *ModelRepository) ListByOwner(ownerID uint) ([]model.ModelAsset, error) {
	var items []model.ModelAsset
	err := r.db.Where("owner_id = ?", ownerID).Order("updated_at desc").Find(&items).Error
	return items, err
}

func (r *ModelRepository) Search(query, tags, robotType string, page, pageSize int) ([]model.ModelAsset, int64, error) {
	return r.ListPublished(query, tags, robotType, "downloads", page, pageSize)
}

func (r *ModelRepository) Top(limit int) ([]model.ModelAsset, error) {
	var items []model.ModelAsset
	err := r.db.Where("status = ?", model.StatusPublished).Order("downloads desc, updated_at desc").Limit(limit).Preload("Owner").Find(&items).Error
	return items, err
}

func (r *ModelRepository) ListPublishedByIDs(ids []uint) ([]model.ModelAsset, error) {
	if len(ids) == 0 {
		return []model.ModelAsset{}, nil
	}
	var items []model.ModelAsset
	err := r.db.Where("status = ? AND id IN ?", model.StatusPublished, ids).Preload("Owner").Find(&items).Error
	return items, err
}

func (r *ModelRepository) CreateVersion(version *model.ModelVersion) error {
	return r.db.Create(version).Error
}

func (r *ModelRepository) IncrementDownloads(id uint) error {
	return r.db.Model(&model.ModelAsset{}).Where("id = ?", id).UpdateColumn("downloads", gorm.Expr("downloads + 1")).Error
}

func (r *ModelRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&model.ModelAsset{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

func (r *ModelRepository) ListPending() ([]model.ModelAsset, error) {
	var items []model.ModelAsset
	err := r.db.Where("status = ?", model.StatusPending).Preload("Owner").Order("updated_at desc").Find(&items).Error
	return items, err
}

type DatasetRepository struct {
	db *gorm.DB
}

func NewDatasetRepository(db *gorm.DB) *DatasetRepository {
	return &DatasetRepository{db: db}
}

func (r *DatasetRepository) Create(item *model.Dataset) error {
	return r.db.Create(item).Error
}

func (r *DatasetRepository) Update(item *model.Dataset) error {
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(item).Error
}

func (r *DatasetRepository) GetByID(id uint) (*model.Dataset, error) {
	var item model.Dataset
	err := r.db.Preload("Owner").Preload("Versions", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at desc")
	}).Preload("Samples").
		First(&item, id).Error
	return &item, err
}

func (r *DatasetRepository) ListPublished(query, tags, sort string, page, pageSize int) ([]model.Dataset, int64, error) {
	db := r.db.Model(&model.Dataset{}).Where("status = ?", model.StatusPublished)
	db = applyTextSearch(db, query, []string{"name", "summary", "description", "scene"})
	db = applyTagFilter(db, tags)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var items []model.Dataset
	err := applyPagination(applySort(db, sort, "downloads desc, updated_at desc"), page, pageSize).
		Preload("Owner").
		Find(&items).Error
	return items, total, err
}

func (r *DatasetRepository) ListByOwner(ownerID uint) ([]model.Dataset, error) {
	var items []model.Dataset
	err := r.db.Where("owner_id = ?", ownerID).Order("updated_at desc").Find(&items).Error
	return items, err
}

func (r *DatasetRepository) Search(query, tags string, page, pageSize int) ([]model.Dataset, int64, error) {
	return r.ListPublished(query, tags, "downloads", page, pageSize)
}

func (r *DatasetRepository) Top(limit int) ([]model.Dataset, error) {
	var items []model.Dataset
	err := r.db.Where("status = ?", model.StatusPublished).Order("downloads desc, updated_at desc").Limit(limit).Preload("Owner").Find(&items).Error
	return items, err
}

func (r *DatasetRepository) ListPublishedByIDs(ids []uint) ([]model.Dataset, error) {
	if len(ids) == 0 {
		return []model.Dataset{}, nil
	}
	var items []model.Dataset
	err := r.db.Where("status = ? AND id IN ?", model.StatusPublished, ids).Preload("Owner").Find(&items).Error
	return items, err
}

func (r *DatasetRepository) CreateVersion(version *model.DatasetVersion) error {
	return r.db.Create(version).Error
}

func (r *DatasetRepository) IncrementDownloads(id uint) error {
	return r.db.Model(&model.Dataset{}).Where("id = ?", id).UpdateColumn("downloads", gorm.Expr("downloads + 1")).Error
}

func (r *DatasetRepository) CountByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Dataset{}).Where("status = ?", status).Count(&count).Error
	return count, err
}

func (r *DatasetRepository) ListPending() ([]model.Dataset, error) {
	var items []model.Dataset
	err := r.db.Where("status = ?", model.StatusPending).Preload("Owner").Order("updated_at desc").Find(&items).Error
	return items, err
}

func (r *DatasetRepository) CreateDownloadTask(item *model.DownloadPackageTask) error {
	return r.db.Create(item).Error
}

func (r *DatasetRepository) ListDownloadTasks(datasetID, userID uint) ([]model.DownloadPackageTask, error) {
	var items []model.DownloadPackageTask
	err := r.db.Where("dataset_id = ? AND user_id = ?", datasetID, userID).Order("created_at desc").Find(&items).Error
	return items, err
}

func (r *DatasetRepository) CreateAccessRequest(item *model.DatasetAccessRequest) error {
	return r.db.Create(item).Error
}

func (r *DatasetRepository) LatestAccessRequest(datasetID, userID uint) (*model.DatasetAccessRequest, error) {
	var item model.DatasetAccessRequest
	err := r.db.Where("dataset_id = ? AND user_id = ?", datasetID, userID).Order("created_at desc").First(&item).Error
	return &item, err
}

func (r *DatasetRepository) FindAccessRequest(id uint) (*model.DatasetAccessRequest, error) {
	var item model.DatasetAccessRequest
	err := r.db.First(&item, id).Error
	return &item, err
}

func (r *DatasetRepository) UpdateAccessRequest(item *model.DatasetAccessRequest) error {
	return r.db.Save(item).Error
}

func (r *DatasetRepository) ListAccessRequests() ([]model.DatasetAccessRequest, error) {
	var items []model.DatasetAccessRequest
	err := r.db.Order("created_at desc").Find(&items).Error
	return items, err
}

type ContentRepository struct {
	db *gorm.DB
}

func NewContentRepository(db *gorm.DB) *ContentRepository {
	return &ContentRepository{db: db}
}

func (r *ContentRepository) ListTemplates() ([]model.TaskTemplate, error) {
	var items []model.TaskTemplate
	err := r.db.Where("status = ?", model.StatusPublished).Order("usage_count desc, updated_at desc").Find(&items).Error
	return items, err
}

func (r *ContentRepository) ListAllTemplates() ([]model.TaskTemplate, error) {
	var items []model.TaskTemplate
	err := r.db.Order("updated_at desc").Find(&items).Error
	return items, err
}

func (r *ContentRepository) ListApplicationCases() ([]model.ApplicationCase, error) {
	var items []model.ApplicationCase
	err := r.db.Where("status = ?", model.StatusPublished).Order("updated_at desc").Find(&items).Error
	return items, err
}

func (r *ContentRepository) ListAllApplicationCases() ([]model.ApplicationCase, error) {
	var items []model.ApplicationCase
	err := r.db.Order("updated_at desc").Find(&items).Error
	return items, err
}

func (r *ContentRepository) GetTemplate(id uint) (*model.TaskTemplate, error) {
	var item model.TaskTemplate
	err := r.db.First(&item, id).Error
	return &item, err
}

func (r *ContentRepository) GetApplicationCase(id uint) (*model.ApplicationCase, error) {
	var item model.ApplicationCase
	err := r.db.First(&item, id).Error
	return &item, err
}

func (r *ContentRepository) ListTemplatesByIDs(ids []uint) ([]model.TaskTemplate, error) {
	if len(ids) == 0 {
		return []model.TaskTemplate{}, nil
	}
	var items []model.TaskTemplate
	err := r.db.Where("status = ? AND id IN ?", model.StatusPublished, ids).Find(&items).Error
	return items, err
}

func (r *ContentRepository) ListApplicationCasesByIDs(ids []uint) ([]model.ApplicationCase, error) {
	if len(ids) == 0 {
		return []model.ApplicationCase{}, nil
	}
	var items []model.ApplicationCase
	err := r.db.Where("status = ? AND id IN ?", model.StatusPublished, ids).Find(&items).Error
	return items, err
}

func (r *ContentRepository) CreateTemplate(item *model.TaskTemplate) error {
	return r.db.Create(item).Error
}

func (r *ContentRepository) UpdateTemplate(item *model.TaskTemplate) error {
	return r.db.Save(item).Error
}

func (r *ContentRepository) DeleteTemplate(id uint) error {
	return r.db.Delete(&model.TaskTemplate{}, id).Error
}

func (r *ContentRepository) BatchUpdateTemplateStatus(ids []uint, status string) error {
	return r.db.Model(&model.TaskTemplate{}).Where("id IN ?", ids).Update("status", status).Error
}

func (r *ContentRepository) BatchDeleteTemplates(ids []uint) error {
	return r.db.Delete(&model.TaskTemplate{}, ids).Error
}

func (r *ContentRepository) CreateApplicationCase(item *model.ApplicationCase) error {
	return r.db.Create(item).Error
}

func (r *ContentRepository) UpdateApplicationCase(item *model.ApplicationCase) error {
	return r.db.Save(item).Error
}

func (r *ContentRepository) DeleteApplicationCase(id uint) error {
	return r.db.Delete(&model.ApplicationCase{}, id).Error
}

func (r *ContentRepository) BatchUpdateApplicationCaseStatus(ids []uint, status string) error {
	return r.db.Model(&model.ApplicationCase{}).Where("id IN ?", ids).Update("status", status).Error
}

func (r *ContentRepository) BatchDeleteApplicationCases(ids []uint) error {
	return r.db.Delete(&model.ApplicationCase{}, ids).Error
}

func (r *ContentRepository) SearchTemplates(query string, page, pageSize int) ([]model.TaskTemplate, int64, error) {
	db := r.db.Model(&model.TaskTemplate{}).Where("status = ?", model.StatusPublished)
	db = applyTextSearch(db, query, []string{"name", "summary", "description", "scene"})

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []model.TaskTemplate
	err := applyPagination(applySort(db, "downloads", "usage_count desc, updated_at desc"), page, pageSize).Find(&items).Error
	return items, total, err
}

func (r *ContentRepository) ListDocCategories(docType string) ([]model.DocumentCategory, error) {
	var items []model.DocumentCategory
	db := r.db.Model(&model.DocumentCategory{})
	if docType != "" {
		db = db.Where("doc_type = ?", docType)
	}
	err := db.Order("name asc").Find(&items).Error
	return items, err
}

func (r *ContentRepository) CreateDocCategory(item *model.DocumentCategory) error {
	return r.db.Create(item).Error
}

func (r *ContentRepository) UpdateDocCategory(item *model.DocumentCategory) error {
	return r.db.Save(item).Error
}

func (r *ContentRepository) GetDocCategory(id uint) (*model.DocumentCategory, error) {
	var item model.DocumentCategory
	err := r.db.First(&item, id).Error
	return &item, err
}

func (r *ContentRepository) DeleteDocCategory(id uint) error {
	return r.db.Delete(&model.DocumentCategory{}, id).Error
}

func (r *ContentRepository) CountDocumentsByCategory(categoryID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Document{}).Where("category_id = ?", categoryID).Count(&count).Error
	return count, err
}

func (r *ContentRepository) ListDocuments(docType string) ([]model.Document, error) {
	var items []model.Document
	db := r.db.Preload("Category").Where("status = ?", model.StatusPublished)
	if docType != "" {
		db = db.Where("doc_type = ?", docType)
	}
	err := db.Order("updated_at desc").Find(&items).Error
	return items, err
}

func (r *ContentRepository) ListAllDocuments() ([]model.Document, error) {
	var items []model.Document
	err := r.db.Preload("Category").Order("updated_at desc").Find(&items).Error
	return items, err
}

func (r *ContentRepository) GetDocument(id uint) (*model.Document, error) {
	var item model.Document
	err := r.db.Preload("Category").First(&item, id).Error
	return &item, err
}

func (r *ContentRepository) CreateDocument(item *model.Document) error {
	return r.db.Create(item).Error
}

func (r *ContentRepository) UpdateDocument(item *model.Document) error {
	return r.db.Save(item).Error
}

func (r *ContentRepository) DeleteDocument(id uint) error {
	return r.db.Delete(&model.Document{}, id).Error
}

func (r *ContentRepository) BatchUpdateDocumentStatus(ids []uint, status string) error {
	return r.db.Model(&model.Document{}).Where("id IN ?", ids).Update("status", status).Error
}

func (r *ContentRepository) BatchDeleteDocuments(ids []uint) error {
	return r.db.Delete(&model.Document{}, ids).Error
}

func (r *ContentRepository) SearchDocuments(query string, page, pageSize int) ([]model.Document, int64, error) {
	db := r.db.Model(&model.Document{}).Where("status = ?", model.StatusPublished)
	db = applyTextSearch(db, query, []string{"title", "summary", "content"})
	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var items []model.Document
	err := applyPagination(db.Order("updated_at desc"), page, pageSize).Preload("Category").Find(&items).Error
	return items, total, err
}

func (r *ContentRepository) ListFAQs() ([]model.FAQ, error) {
	var items []model.FAQ
	err := r.db.Order("updated_at desc").Find(&items).Error
	return items, err
}

func (r *ContentRepository) CreateFAQ(item *model.FAQ) error {
	return r.db.Create(item).Error
}

func (r *ContentRepository) UpdateFAQ(item *model.FAQ) error {
	return r.db.Save(item).Error
}

func (r *ContentRepository) DeleteFAQ(id uint) error {
	return r.db.Delete(&model.FAQ{}, id).Error
}

func (r *ContentRepository) GetFAQ(id uint) (*model.FAQ, error) {
	var item model.FAQ
	err := r.db.First(&item, id).Error
	return &item, err
}

func (r *ContentRepository) ListVideoTutorials(activeOnly bool) ([]model.VideoTutorial, error) {
	var items []model.VideoTutorial
	db := r.db.Model(&model.VideoTutorial{})
	if activeOnly {
		db = db.Where("active = ?", true)
	}
	err := db.Order("sort_order asc, updated_at desc").Find(&items).Error
	return items, err
}

func (r *ContentRepository) ListAgreementTemplates(activeOnly bool) ([]model.AgreementTemplate, error) {
	var items []model.AgreementTemplate
	db := r.db.Model(&model.AgreementTemplate{})
	if activeOnly {
		db = db.Where("active = ?", true)
	}
	err := db.Order("sort_order asc, updated_at desc").Find(&items).Error
	return items, err
}

func (r *ContentRepository) GetAgreementTemplate(id uint) (*model.AgreementTemplate, error) {
	var item model.AgreementTemplate
	err := r.db.First(&item, id).Error
	return &item, err
}

func (r *ContentRepository) CreateAgreementTemplate(item *model.AgreementTemplate) error {
	return r.db.Create(item).Error
}

func (r *ContentRepository) UpdateAgreementTemplate(item *model.AgreementTemplate) error {
	return r.db.Save(item).Error
}

func (r *ContentRepository) DeleteAgreementTemplate(id uint) error {
	return r.db.Delete(&model.AgreementTemplate{}, id).Error
}

func (r *ContentRepository) ListDatasetPrivacyOptions(activeOnly bool) ([]model.DatasetPrivacyOption, error) {
	var items []model.DatasetPrivacyOption
	db := r.db.Model(&model.DatasetPrivacyOption{})
	if activeOnly {
		db = db.Where("active = ?", true)
	}
	err := db.Order("sort_order asc, updated_at desc").Find(&items).Error
	return items, err
}

func (r *ContentRepository) GetDatasetPrivacyOption(id uint) (*model.DatasetPrivacyOption, error) {
	var item model.DatasetPrivacyOption
	err := r.db.First(&item, id).Error
	return &item, err
}

func (r *ContentRepository) CreateDatasetPrivacyOption(item *model.DatasetPrivacyOption) error {
	return r.db.Create(item).Error
}

func (r *ContentRepository) UpdateDatasetPrivacyOption(item *model.DatasetPrivacyOption) error {
	return r.db.Save(item).Error
}

func (r *ContentRepository) DeleteDatasetPrivacyOption(id uint) error {
	return r.db.Delete(&model.DatasetPrivacyOption{}, id).Error
}

func (r *ContentRepository) GetVideoTutorial(id uint) (*model.VideoTutorial, error) {
	var item model.VideoTutorial
	err := r.db.First(&item, id).Error
	return &item, err
}

func (r *ContentRepository) CreateVideoTutorial(item *model.VideoTutorial) error {
	return r.db.Create(item).Error
}

func (r *ContentRepository) UpdateVideoTutorial(item *model.VideoTutorial) error {
	return r.db.Save(item).Error
}

func (r *ContentRepository) DeleteVideoTutorial(id uint) error {
	return r.db.Delete(&model.VideoTutorial{}, id).Error
}

func (r *ContentRepository) BatchUpdateVideoTutorialStatus(ids []uint, active bool) error {
	return r.db.Model(&model.VideoTutorial{}).Where("id IN ?", ids).Update("active", active).Error
}

func (r *ContentRepository) BatchDeleteVideoTutorials(ids []uint) error {
	return r.db.Delete(&model.VideoTutorial{}, ids).Error
}

type UserActivityRepository struct {
	db *gorm.DB
}

func NewUserActivityRepository(db *gorm.DB) *UserActivityRepository {
	return &UserActivityRepository{db: db}
}

func (r *UserActivityRepository) AddFavorite(item *model.Favorite) error {
	return r.db.Where("user_id = ? and resource_type = ? and resource_id = ?", item.UserID, item.ResourceType, item.ResourceID).FirstOrCreate(item).Error
}

func (r *UserActivityRepository) RemoveFavorite(userID uint, resourceType string, resourceID uint) error {
	return r.db.Where("user_id = ? and resource_type = ? and resource_id = ?", userID, resourceType, resourceID).Delete(&model.Favorite{}).Error
}

func (r *UserActivityRepository) IsFavorited(userID uint, resourceType string, resourceID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Favorite{}).Where("user_id = ? and resource_type = ? and resource_id = ?", userID, resourceType, resourceID).Count(&count).Error
	return count > 0, err
}

func (r *UserActivityRepository) ListFavorites(userID uint) ([]model.Favorite, error) {
	var items []model.Favorite
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&items).Error
	return items, err
}

func (r *UserActivityRepository) AddDownload(record *model.DownloadRecord) error {
	return r.db.Create(record).Error
}

func (r *UserActivityRepository) ListDownloads(userID uint) ([]model.DownloadRecord, error) {
	var items []model.DownloadRecord
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&items).Error
	return items, err
}

func (r *UserActivityRepository) AddNotification(item *model.Notification) error {
	return r.db.Create(item).Error
}

func (r *UserActivityRepository) ListNotifications(userID uint) ([]model.Notification, error) {
	var items []model.Notification
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&items).Error
	return items, err
}

func (r *UserActivityRepository) MarkNotificationsRead(userID uint) error {
	return r.db.Model(&model.Notification{}).Where("user_id = ?", userID).Update("read", true).Error
}

func (r *UserActivityRepository) ConfirmAgreement(item *model.AgreementRecord) error {
	return r.db.Where("user_id = ? and dataset_id = ?", item.UserID, item.DatasetID).FirstOrCreate(item).Error
}

func (r *UserActivityRepository) HasAgreement(userID, datasetID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.AgreementRecord{}).Where("user_id = ? and dataset_id = ?", userID, datasetID).Count(&count).Error
	return count > 0, err
}

type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) AddLog(item *model.ReviewLog) error {
	return r.db.Create(item).Error
}

type OperationLogRepository struct {
	db *gorm.DB
}

func NewOperationLogRepository(db *gorm.DB) *OperationLogRepository {
	return &OperationLogRepository{db: db}
}

func (r *OperationLogRepository) Create(item *model.AdminOperationLog) error {
	return r.db.Create(item).Error
}

func (r *OperationLogRepository) ListRecent(limit int) ([]model.AdminOperationLog, error) {
	var items []model.AdminOperationLog
	err := r.db.Order("created_at desc").Limit(limit).Find(&items).Error
	return items, err
}

func applyTextSearch(db *gorm.DB, query string, columns []string) *gorm.DB {
	if strings.TrimSpace(query) == "" {
		return db
	}

	like := "%" + strings.TrimSpace(query) + "%"
	conditions := make([]string, 0, len(columns))
	args := make([]any, 0, len(columns))
	for _, column := range columns {
		conditions = append(conditions, fmt.Sprintf("%s LIKE ?", column))
		args = append(args, like)
	}
	return db.Where(strings.Join(conditions, " OR "), args...)
}

func applyTagFilter(db *gorm.DB, tags string) *gorm.DB {
	for _, tag := range strings.Split(tags, ",") {
		trimmed := strings.TrimSpace(tag)
		if trimmed != "" {
			db = db.Where("tags LIKE ?", "%"+trimmed+"%")
		}
	}
	return db
}

func applySort(db *gorm.DB, sort, fallback string) *gorm.DB {
	switch sort {
	case "latest":
		return db.Order("updated_at desc")
	case "name":
		return db.Order("name asc")
	case "downloads", "hot":
		return db.Order("downloads desc, updated_at desc")
	default:
		return db.Order(fallback)
	}
}

func applyPagination(db *gorm.DB, page, pageSize int) *gorm.DB {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 12
	}
	return db.Offset((page - 1) * pageSize).Limit(pageSize)
}

func EnsureSeeded(db *gorm.DB, seedFunc func(tx *gorm.DB) error) error {
	var count int64
	if err := db.Model(&model.User{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return db.Transaction(seedFunc)
}
