package service

import (
	"net/http"
	"strings"
	"time"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/model"
	"github.com/Andras5014/EIBotHub/server/internal/repository"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

type AdminService struct {
	users    *repository.UserRepository
	portal   *repository.PortalRepository
	models   *repository.ModelRepository
	datasets *repository.DatasetRepository
	content  *repository.ContentRepository
	reviews  *repository.ReviewRepository
	ops      *repository.OperationLogRepository
}

func NewAdminService(users *repository.UserRepository, portal *repository.PortalRepository, models *repository.ModelRepository, datasets *repository.DatasetRepository, content *repository.ContentRepository, reviews *repository.ReviewRepository, ops *repository.OperationLogRepository) *AdminService {
	return &AdminService{users: users, portal: portal, models: models, datasets: datasets, content: content, reviews: reviews, ops: ops}
}

func (s *AdminService) Dashboard() (*dto.DashboardResponse, error) {
	users, err := s.users.Count()
	if err != nil {
		return nil, err
	}
	publishedModels, err := s.models.CountByStatus(model.StatusPublished)
	if err != nil {
		return nil, err
	}
	pendingModels, err := s.models.CountByStatus(model.StatusPending)
	if err != nil {
		return nil, err
	}
	publishedDatasets, err := s.datasets.CountByStatus(model.StatusPublished)
	if err != nil {
		return nil, err
	}
	pendingDatasets, err := s.datasets.CountByStatus(model.StatusPending)
	if err != nil {
		return nil, err
	}
	announcements, err := s.portal.CountAnnouncements()
	if err != nil {
		return nil, err
	}
	return &dto.DashboardResponse{
		Users:             users,
		PublishedModels:   publishedModels,
		PendingModels:     pendingModels,
		PublishedDatasets: publishedDatasets,
		PendingDatasets:   pendingDatasets,
		Announcements:     announcements,
	}, nil
}

func (s *AdminService) Reviews(resourceType string) ([]dto.ReviewItem, error) {
	switch resourceType {
	case "models":
		items, err := s.models.ListPending()
		if err != nil {
			return nil, err
		}
		result := make([]dto.ReviewItem, 0, len(items))
		for _, item := range items {
			result = append(result, dto.ReviewItem{
				ID:        item.ID,
				Type:      "models",
				Title:     item.Name,
				Summary:   item.Summary,
				Status:    item.Status,
				Owner:     item.Owner.Username,
				UpdatedAt: item.UpdatedAt,
			})
		}
		return result, nil
	case "datasets":
		items, err := s.datasets.ListPending()
		if err != nil {
			return nil, err
		}
		result := make([]dto.ReviewItem, 0, len(items))
		for _, item := range items {
			result = append(result, dto.ReviewItem{
				ID:        item.ID,
				Type:      "datasets",
				Title:     item.Name,
				Summary:   item.Summary,
				Status:    item.Status,
				Owner:     item.Owner.Username,
				UpdatedAt: item.UpdatedAt,
			})
		}
		return result, nil
	default:
		return nil, support.NewError(http.StatusBadRequest, "invalid_review_type", "unsupported review type")
	}
}

func (s *AdminService) Decide(resourceType string, resourceID, reviewerID uint, input dto.ReviewDecisionRequest) error {
	var err error
	switch resourceType {
	case "models":
		item, getErr := s.models.GetByID(resourceID)
		if getErr != nil {
			return getErr
		}
		item.Status = mapDecision(input.Decision)
		err = s.models.Update(item)
	case "datasets":
		item, getErr := s.datasets.GetByID(resourceID)
		if getErr != nil {
			return getErr
		}
		item.Status = mapDecision(input.Decision)
		err = s.datasets.Update(item)
	default:
		return support.NewError(http.StatusBadRequest, "invalid_review_type", "unsupported review type")
	}
	if err != nil {
		return err
	}
	return s.reviews.AddLog(&model.ReviewLog{
		ResourceType: resourceType,
		ResourceID:   resourceID,
		ReviewerID:   reviewerID,
		Decision:     input.Decision,
		Comment:      input.Comment,
		CreatedAt:    time.Now(),
	})
}

func (s *AdminService) Announcements() ([]dto.AnnouncementItem, error) {
	items, err := s.portal.ListAnnouncements()
	if err != nil {
		return nil, err
	}
	result := make([]dto.AnnouncementItem, 0, len(items))
	for _, item := range items {
		result = append(result, dto.AnnouncementItem{
			ID:          item.ID,
			Title:       item.Title,
			Summary:     item.Summary,
			Link:        item.Link,
			Pinned:      item.Pinned,
			PublishedAt: item.PublishedAt,
		})
	}
	return result, nil
}

func (s *AdminService) CreateAnnouncement(input dto.AnnouncementRequest) (*dto.AnnouncementItem, error) {
	item := &model.Announcement{
		Title:       input.Title,
		Summary:     input.Summary,
		Link:        input.Link,
		Pinned:      input.Pinned,
		PublishedAt: time.Now(),
	}
	if err := s.portal.CreateAnnouncement(item); err != nil {
		return nil, err
	}
	return &dto.AnnouncementItem{
		ID:          item.ID,
		Title:       item.Title,
		Summary:     item.Summary,
		Link:        item.Link,
		Pinned:      item.Pinned,
		PublishedAt: item.PublishedAt,
	}, nil
}

func (s *AdminService) DeleteTemplate(id, adminUserID uint) error {
	item, err := s.content.GetTemplate(id)
	if err != nil {
		return err
	}
	if err := s.content.DeleteTemplate(id); err != nil {
		return err
	}
	return s.logOperation(adminUserID, "delete", "task-template", id, item.Name, item.Summary)
}

func (s *AdminService) BatchUpdateTemplateStatus(adminUserID uint, input dto.BatchStatusRequest) error {
	if err := s.content.BatchUpdateTemplateStatus(input.IDs, input.Status); err != nil {
		return err
	}
	return s.logOperation(adminUserID, "batch_status", "task-template", 0, "批量更新模板状态", input.Status)
}

func (s *AdminService) BatchDeleteTemplates(adminUserID uint, input dto.BatchDeleteRequest) error {
	if err := s.content.BatchDeleteTemplates(input.IDs); err != nil {
		return err
	}
	return s.logOperation(adminUserID, "batch_delete", "task-template", 0, "批量删除模板", joinUintIDs(input.IDs))
}

func (s *AdminService) ModuleSettings() ([]dto.ModuleSettingItem, error) {
	items, err := s.portal.ListModuleSettings()
	if err != nil {
		return nil, err
	}
	byKey := make(map[string]model.HomeModuleSetting, len(items))
	for _, item := range items {
		byKey[item.ModuleKey] = item
	}
	result := make([]dto.ModuleSettingItem, 0, len(defaultHomeModuleKeys))
	for _, key := range defaultHomeModuleKeys {
		if item, ok := byKey[key]; ok {
			result = append(result, dto.ModuleSettingItem{
				ID:        item.ID,
				ModuleKey: key,
				Label:     homeModuleLabels[key],
				Enabled:   item.Enabled,
				UpdatedAt: item.UpdatedAt,
			})
			continue
		}
		result = append(result, dto.ModuleSettingItem{
			ModuleKey: key,
			Label:     homeModuleLabels[key],
			Enabled:   true,
		})
	}
	return result, nil
}

func (s *AdminService) UpdateModuleSetting(moduleKey string, input dto.ModuleSettingRequest) error {
	if _, ok := homeModuleLabels[moduleKey]; !ok {
		return support.NewError(http.StatusBadRequest, "invalid_module_key", "unsupported home module")
	}
	return s.portal.UpsertModuleSetting(&model.HomeModuleSetting{
		ModuleKey: moduleKey,
		Enabled:   input.Enabled,
	})
}

func (s *AdminService) FeaturedResources() ([]dto.FeaturedResourceItem, error) {
	items, err := s.portal.ListFeaturedResources()
	if err != nil {
		return nil, err
	}
	result := make([]dto.FeaturedResourceItem, 0, len(items))
	for _, item := range items {
		title, summary, route := s.resolveFeaturedResource(item)
		result = append(result, dto.FeaturedResourceItem{
			ID:           item.ID,
			ResourceType: item.ResourceType,
			ResourceID:   item.ResourceID,
			Title:        title,
			Summary:      summary,
			Route:        route,
			SortOrder:    item.SortOrder,
			Enabled:      item.Enabled,
			UpdatedAt:    item.UpdatedAt,
		})
	}
	return result, nil
}

func (s *AdminService) CreateFeaturedResource(input dto.FeaturedResourceRequest) (*dto.FeaturedResourceItem, error) {
	item := &model.FeaturedResource{
		ResourceType: input.ResourceType,
		ResourceID:   input.ResourceID,
		SortOrder:    input.SortOrder,
		Enabled:      input.Enabled,
	}
	if err := s.portal.CreateFeaturedResource(item); err != nil {
		return nil, err
	}
	title, summary, route := s.resolveFeaturedResource(*item)
	return &dto.FeaturedResourceItem{
		ID:           item.ID,
		ResourceType: item.ResourceType,
		ResourceID:   item.ResourceID,
		Title:        title,
		Summary:      summary,
		Route:        route,
		SortOrder:    item.SortOrder,
		Enabled:      item.Enabled,
		UpdatedAt:    item.UpdatedAt,
	}, nil
}

func (s *AdminService) UpdateFeaturedResource(id uint, input dto.FeaturedResourceRequest) (*dto.FeaturedResourceItem, error) {
	item, err := s.portal.GetFeaturedResource(id)
	if err != nil {
		return nil, err
	}
	item.ResourceType = input.ResourceType
	item.ResourceID = input.ResourceID
	item.SortOrder = input.SortOrder
	item.Enabled = input.Enabled
	if err := s.portal.UpdateFeaturedResource(item); err != nil {
		return nil, err
	}
	title, summary, route := s.resolveFeaturedResource(*item)
	result := &dto.FeaturedResourceItem{
		ID:           item.ID,
		ResourceType: item.ResourceType,
		ResourceID:   item.ResourceID,
		Title:        title,
		Summary:      summary,
		Route:        route,
		SortOrder:    item.SortOrder,
		Enabled:      item.Enabled,
		UpdatedAt:    item.UpdatedAt,
	}
	return result, nil
}

func (s *AdminService) DeleteFeaturedResource(id uint) error {
	return s.portal.DeleteFeaturedResource(id)
}

func (s *AdminService) Templates() ([]dto.TaskTemplateItem, error) {
	items, err := s.content.ListAllTemplates()
	if err != nil {
		return nil, err
	}
	result := make([]dto.TaskTemplateItem, 0, len(items))
	for _, item := range items {
		result = append(result, toTemplateItem(item))
	}
	return result, nil
}

func (s *AdminService) CreateTemplate(input dto.TaskTemplateAdminRequest) (*dto.TaskTemplateItem, error) {
	item := &model.TaskTemplate{
		Name:        input.Name,
		Summary:     input.Summary,
		Description: input.Description,
		Category:    input.Category,
		Scene:       input.Scene,
		Guide:       input.Guide,
		ResourceRef: input.ResourceRef,
		Status:      input.Status,
	}
	if err := s.content.CreateTemplate(item); err != nil {
		return nil, err
	}
	result := toTemplateItem(*item)
	return &result, nil
}

func (s *AdminService) UpdateTemplate(id uint, input dto.TaskTemplateAdminRequest) (*dto.TaskTemplateItem, error) {
	item, err := s.content.GetTemplate(id)
	if err != nil {
		return nil, err
	}
	item.Name = input.Name
	item.Summary = input.Summary
	item.Description = input.Description
	item.Category = input.Category
	item.Scene = input.Scene
	item.Guide = input.Guide
	item.ResourceRef = input.ResourceRef
	item.Status = input.Status
	if err := s.content.UpdateTemplate(item); err != nil {
		return nil, err
	}
	result := toTemplateItem(*item)
	return &result, nil
}

func (s *AdminService) ApplicationCases() ([]dto.ApplicationCaseItem, error) {
	items, err := s.content.ListAllApplicationCases()
	if err != nil {
		return nil, err
	}
	result := make([]dto.ApplicationCaseItem, 0, len(items))
	for _, item := range items {
		result = append(result, toApplicationCase(item))
	}
	return result, nil
}

func (s *AdminService) CreateApplicationCase(input dto.ApplicationCaseAdminRequest) (*dto.ApplicationCaseItem, error) {
	item := &model.ApplicationCase{
		Title:      input.Title,
		Summary:    input.Summary,
		Category:   input.Category,
		Guide:      input.Guide,
		CoverImage: input.CoverImage,
		Status:     input.Status,
	}
	if err := s.content.CreateApplicationCase(item); err != nil {
		return nil, err
	}
	result := toApplicationCase(*item)
	return &result, nil
}

func (s *AdminService) DeleteApplicationCase(id, adminUserID uint) error {
	item, err := s.content.GetApplicationCase(id)
	if err != nil {
		return err
	}
	if err := s.content.DeleteApplicationCase(id); err != nil {
		return err
	}
	return s.logOperation(adminUserID, "delete", "application-case", id, item.Title, item.Summary)
}

func (s *AdminService) BatchUpdateApplicationCaseStatus(adminUserID uint, input dto.BatchStatusRequest) error {
	if err := s.content.BatchUpdateApplicationCaseStatus(input.IDs, input.Status); err != nil {
		return err
	}
	return s.logOperation(adminUserID, "batch_status", "application-case", 0, "批量更新案例状态", input.Status)
}

func (s *AdminService) BatchDeleteApplicationCases(adminUserID uint, input dto.BatchDeleteRequest) error {
	if err := s.content.BatchDeleteApplicationCases(input.IDs); err != nil {
		return err
	}
	return s.logOperation(adminUserID, "batch_delete", "application-case", 0, "批量删除案例", joinUintIDs(input.IDs))
}

func (s *AdminService) UpdateApplicationCase(id uint, input dto.ApplicationCaseAdminRequest) (*dto.ApplicationCaseItem, error) {
	item, err := s.content.GetApplicationCase(id)
	if err != nil {
		return nil, err
	}
	item.Title = input.Title
	item.Summary = input.Summary
	item.Category = input.Category
	item.Guide = input.Guide
	item.CoverImage = input.CoverImage
	item.Status = input.Status
	if err := s.content.UpdateApplicationCase(item); err != nil {
		return nil, err
	}
	result := toApplicationCase(*item)
	return &result, nil
}

func (s *AdminService) DocCategories() ([]dto.DocumentCategoryItem, error) {
	items, err := s.content.ListDocCategories("")
	if err != nil {
		return nil, err
	}
	result := make([]dto.DocumentCategoryItem, 0, len(items))
	for _, item := range items {
		result = append(result, dto.DocumentCategoryItem{
			ID:      item.ID,
			Name:    item.Name,
			DocType: item.DocType,
		})
	}
	return result, nil
}

func (s *AdminService) CreateDocCategory(input dto.DocumentCategoryRequest) (*dto.DocumentCategoryItem, error) {
	item := &model.DocumentCategory{
		Name:    input.Name,
		DocType: input.DocType,
	}
	if err := s.content.CreateDocCategory(item); err != nil {
		return nil, err
	}
	return &dto.DocumentCategoryItem{
		ID:      item.ID,
		Name:    item.Name,
		DocType: item.DocType,
	}, nil
}

func (s *AdminService) UpdateDocCategory(id uint, input dto.DocumentCategoryRequest) (*dto.DocumentCategoryItem, error) {
	item, err := s.content.GetDocCategory(id)
	if err != nil {
		return nil, err
	}
	item.Name = input.Name
	item.DocType = input.DocType
	if err := s.content.UpdateDocCategory(item); err != nil {
		return nil, err
	}
	return &dto.DocumentCategoryItem{
		ID:      item.ID,
		Name:    item.Name,
		DocType: item.DocType,
	}, nil
}

func (s *AdminService) DeleteDocCategory(id, adminUserID uint) error {
	count, err := s.content.CountDocumentsByCategory(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return support.NewError(http.StatusBadRequest, "category_has_documents", "当前分类下仍有文档，不能直接删除")
	}
	item, err := s.content.GetDocCategory(id)
	if err != nil {
		return err
	}
	if err := s.content.DeleteDocCategory(id); err != nil {
		return err
	}
	return s.logOperation(adminUserID, "delete", "doc-category", id, item.Name, item.DocType)
}

func (s *AdminService) Documents() ([]dto.DocumentItem, error) {
	items, err := s.content.ListAllDocuments()
	if err != nil {
		return nil, err
	}
	result := make([]dto.DocumentItem, 0, len(items))
	for _, item := range items {
		result = append(result, toDocumentItem(item))
	}
	return result, nil
}

func (s *AdminService) CreateDocument(input dto.DocumentAdminRequest) (*dto.DocumentItem, error) {
	item := &model.Document{
		CategoryID: input.CategoryID,
		Title:      input.Title,
		Summary:    input.Summary,
		Content:    input.Content,
		DocType:    input.DocType,
		Status:     input.Status,
	}
	if err := s.content.CreateDocument(item); err != nil {
		return nil, err
	}
	created, err := s.content.GetDocument(item.ID)
	if err != nil {
		return nil, err
	}
	result := toDocumentItem(*created)
	return &result, nil
}

func (s *AdminService) DeleteDocument(id, adminUserID uint) error {
	item, err := s.content.GetDocument(id)
	if err != nil {
		return err
	}
	if err := s.content.DeleteDocument(id); err != nil {
		return err
	}
	return s.logOperation(adminUserID, "delete", "doc", id, item.Title, item.Summary)
}

func (s *AdminService) BatchUpdateDocumentStatus(adminUserID uint, input dto.BatchStatusRequest) error {
	if err := s.content.BatchUpdateDocumentStatus(input.IDs, input.Status); err != nil {
		return err
	}
	return s.logOperation(adminUserID, "batch_status", "doc", 0, "批量更新文档状态", input.Status)
}

func (s *AdminService) BatchDeleteDocuments(adminUserID uint, input dto.BatchDeleteRequest) error {
	if err := s.content.BatchDeleteDocuments(input.IDs); err != nil {
		return err
	}
	return s.logOperation(adminUserID, "batch_delete", "doc", 0, "批量删除文档", joinUintIDs(input.IDs))
}

func (s *AdminService) UpdateDocument(id uint, input dto.DocumentAdminRequest) (*dto.DocumentItem, error) {
	item, err := s.content.GetDocument(id)
	if err != nil {
		return nil, err
	}
	item.CategoryID = input.CategoryID
	item.Title = input.Title
	item.Summary = input.Summary
	item.Content = input.Content
	item.DocType = input.DocType
	item.Status = input.Status
	if err := s.content.UpdateDocument(item); err != nil {
		return nil, err
	}
	updated, err := s.content.GetDocument(item.ID)
	if err != nil {
		return nil, err
	}
	result := toDocumentItem(*updated)
	return &result, nil
}

func (s *AdminService) FAQs() ([]dto.FAQItem, error) {
	items, err := s.content.ListFAQs()
	if err != nil {
		return nil, err
	}
	result := make([]dto.FAQItem, 0, len(items))
	for _, item := range items {
		result = append(result, dto.FAQItem{
			ID:        item.ID,
			Question:  item.Question,
			Answer:    item.Answer,
			UpdatedAt: item.UpdatedAt,
		})
	}
	return result, nil
}

func (s *AdminService) CreateFAQ(input dto.FAQRequest) (*dto.FAQItem, error) {
	item := &model.FAQ{
		Question: input.Question,
		Answer:   input.Answer,
	}
	if err := s.content.CreateFAQ(item); err != nil {
		return nil, err
	}
	return &dto.FAQItem{
		ID:        item.ID,
		Question:  item.Question,
		Answer:    item.Answer,
		UpdatedAt: item.UpdatedAt,
	}, nil
}

func (s *AdminService) DeleteFAQ(id, adminUserID uint) error {
	item, err := s.content.GetFAQ(id)
	if err != nil {
		return err
	}
	if err := s.content.DeleteFAQ(id); err != nil {
		return err
	}
	return s.logOperation(adminUserID, "delete", "faq", id, item.Question, item.Answer)
}

func (s *AdminService) UpdateFAQ(id uint, input dto.FAQRequest) (*dto.FAQItem, error) {
	item, err := s.content.GetFAQ(id)
	if err != nil {
		return nil, err
	}
	item.Question = input.Question
	item.Answer = input.Answer
	if err := s.content.UpdateFAQ(item); err != nil {
		return nil, err
	}
	return &dto.FAQItem{
		ID:        item.ID,
		Question:  item.Question,
		Answer:    item.Answer,
		UpdatedAt: item.UpdatedAt,
	}, nil
}

func (s *AdminService) VideoTutorials() ([]dto.VideoTutorialItem, error) {
	items, err := s.content.ListVideoTutorials(false)
	if err != nil {
		return nil, err
	}
	result := make([]dto.VideoTutorialItem, 0, len(items))
	for _, item := range items {
		result = append(result, toVideoTutorialItem(item))
	}
	return result, nil
}

func (s *AdminService) CreateVideoTutorial(input dto.VideoTutorialRequest) (*dto.VideoTutorialItem, error) {
	item := &model.VideoTutorial{
		Title:     input.Title,
		Summary:   input.Summary,
		Link:      input.Link,
		Category:  input.Category,
		SortOrder: input.SortOrder,
		Active:    input.Active,
	}
	if err := s.content.CreateVideoTutorial(item); err != nil {
		return nil, err
	}
	result := toVideoTutorialItem(*item)
	return &result, nil
}

func (s *AdminService) DeleteVideoTutorial(id, adminUserID uint) error {
	item, err := s.content.GetVideoTutorial(id)
	if err != nil {
		return err
	}
	if err := s.content.DeleteVideoTutorial(id); err != nil {
		return err
	}
	return s.logOperation(adminUserID, "delete", "video", id, item.Title, item.Summary)
}

func (s *AdminService) BatchUpdateVideoStatus(adminUserID uint, ids []uint, active bool) error {
	if err := s.content.BatchUpdateVideoTutorialStatus(ids, active); err != nil {
		return err
	}
	state := "disabled"
	if active {
		state = "enabled"
	}
	return s.logOperation(adminUserID, "batch_status", "video", 0, "批量更新视频状态", state)
}

func (s *AdminService) BatchDeleteVideos(adminUserID uint, input dto.BatchDeleteRequest) error {
	if err := s.content.BatchDeleteVideoTutorials(input.IDs); err != nil {
		return err
	}
	return s.logOperation(adminUserID, "batch_delete", "video", 0, "批量删除视频", joinUintIDs(input.IDs))
}

func (s *AdminService) UpdateVideoTutorial(id uint, input dto.VideoTutorialRequest) (*dto.VideoTutorialItem, error) {
	item, err := s.content.GetVideoTutorial(id)
	if err != nil {
		return nil, err
	}
	item.Title = input.Title
	item.Summary = input.Summary
	item.Link = input.Link
	item.Category = input.Category
	item.SortOrder = input.SortOrder
	item.Active = input.Active
	if err := s.content.UpdateVideoTutorial(item); err != nil {
		return nil, err
	}
	result := toVideoTutorialItem(*item)
	return &result, nil
}

func (s *AdminService) PublicVideoTutorials() ([]dto.VideoTutorialItem, error) {
	items, err := s.content.ListVideoTutorials(true)
	if err != nil {
		return nil, err
	}
	result := make([]dto.VideoTutorialItem, 0, len(items))
	for _, item := range items {
		result = append(result, toVideoTutorialItem(item))
	}
	return result, nil
}

func (s *AdminService) AgreementTemplates() ([]dto.AgreementTemplateItem, error) {
	items, err := s.content.ListAgreementTemplates(false)
	if err != nil {
		return nil, err
	}
	result := make([]dto.AgreementTemplateItem, 0, len(items))
	for _, item := range items {
		result = append(result, toAgreementTemplateItem(item))
	}
	return result, nil
}

func (s *AdminService) CreateAgreementTemplate(input dto.AgreementTemplateRequest) (*dto.AgreementTemplateItem, error) {
	item := &model.AgreementTemplate{
		Name:      input.Name,
		Content:   input.Content,
		SortOrder: input.SortOrder,
		Active:    input.Active,
	}
	if err := s.content.CreateAgreementTemplate(item); err != nil {
		return nil, err
	}
	result := toAgreementTemplateItem(*item)
	return &result, nil
}

func (s *AdminService) UpdateAgreementTemplate(id uint, input dto.AgreementTemplateRequest) (*dto.AgreementTemplateItem, error) {
	item, err := s.content.GetAgreementTemplate(id)
	if err != nil {
		return nil, err
	}
	item.Name = input.Name
	item.Content = input.Content
	item.SortOrder = input.SortOrder
	item.Active = input.Active
	if err := s.content.UpdateAgreementTemplate(item); err != nil {
		return nil, err
	}
	result := toAgreementTemplateItem(*item)
	return &result, nil
}

func (s *AdminService) DeleteAgreementTemplate(id uint) error {
	return s.content.DeleteAgreementTemplate(id)
}

func (s *AdminService) PrivacyOptions() ([]dto.DatasetPrivacyOptionItem, error) {
	items, err := s.content.ListDatasetPrivacyOptions(false)
	if err != nil {
		return nil, err
	}
	result := make([]dto.DatasetPrivacyOptionItem, 0, len(items))
	for _, item := range items {
		result = append(result, toDatasetPrivacyOptionItem(item))
	}
	return result, nil
}

func (s *AdminService) CreatePrivacyOption(input dto.DatasetPrivacyOptionRequest) (*dto.DatasetPrivacyOptionItem, error) {
	item := &model.DatasetPrivacyOption{
		Code:        input.Code,
		Name:        input.Name,
		Description: input.Description,
		SortOrder:   input.SortOrder,
		Active:      input.Active,
	}
	if err := s.content.CreateDatasetPrivacyOption(item); err != nil {
		return nil, err
	}
	result := toDatasetPrivacyOptionItem(*item)
	return &result, nil
}

func (s *AdminService) UpdatePrivacyOption(id uint, input dto.DatasetPrivacyOptionRequest) (*dto.DatasetPrivacyOptionItem, error) {
	item, err := s.content.GetDatasetPrivacyOption(id)
	if err != nil {
		return nil, err
	}
	item.Code = input.Code
	item.Name = input.Name
	item.Description = input.Description
	item.SortOrder = input.SortOrder
	item.Active = input.Active
	if err := s.content.UpdateDatasetPrivacyOption(item); err != nil {
		return nil, err
	}
	result := toDatasetPrivacyOptionItem(*item)
	return &result, nil
}

func (s *AdminService) DeletePrivacyOption(id uint) error {
	return s.content.DeleteDatasetPrivacyOption(id)
}

func (s *AdminService) OperationLogs() ([]dto.AdminOperationLogItem, error) {
	items, err := s.ops.ListRecent(100)
	if err != nil {
		return nil, err
	}
	result := make([]dto.AdminOperationLogItem, 0, len(items))
	for _, item := range items {
		adminName := ""
		if admin, adminErr := s.users.FindByID(item.AdminUserID); adminErr == nil {
			adminName = admin.Username
		}
		result = append(result, dto.AdminOperationLogItem{
			ID:           item.ID,
			AdminUserID:  item.AdminUserID,
			AdminName:    adminName,
			Action:       item.Action,
			ResourceType: item.ResourceType,
			ResourceID:   item.ResourceID,
			Summary:      item.Summary,
			Detail:       item.Detail,
			CreatedAt:    item.CreatedAt,
		})
	}
	return result, nil
}

func (s *AdminService) resolveFeaturedResource(item model.FeaturedResource) (string, string, string) {
	switch item.ResourceType {
	case "model":
		resource, err := s.models.GetByID(item.ResourceID)
		if err == nil {
			return resource.Name, resource.Summary, "/models/" + uintToString(resource.ID)
		}
	case "dataset":
		resource, err := s.datasets.GetByID(item.ResourceID)
		if err == nil {
			return resource.Name, resource.Summary, "/datasets/" + uintToString(resource.ID)
		}
	case "task-template":
		resource, err := s.content.GetTemplate(item.ResourceID)
		if err == nil {
			return resource.Name, resource.Summary, "/templates/" + uintToString(resource.ID)
		}
	case "application-case":
		resource, err := s.content.GetApplicationCase(item.ResourceID)
		if err == nil {
			return resource.Title, resource.Summary, "/applications/" + uintToString(resource.ID)
		}
	}
	return "资源不存在", "请检查推荐位配置", ""
}

func mapDecision(decision string) string {
	if decision == "approved" {
		return model.StatusPublished
	}
	return model.StatusRejected
}

func (s *AdminService) logOperation(adminUserID uint, action, resourceType string, resourceID uint, summary, detail string) error {
	if s.ops == nil {
		return nil
	}
	return s.ops.Create(&model.AdminOperationLog{
		AdminUserID:  adminUserID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Summary:      summary,
		Detail:       detail,
		CreatedAt:    time.Now(),
	})
}

func joinUintIDs(ids []uint) string {
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, uintToString(id))
	}
	return strings.Join(parts, ",")
}
