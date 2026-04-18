package service

import (
	"net/http"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/repository"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

type CatalogService struct {
	content *repository.ContentRepository
}

func NewCatalogService(content *repository.ContentRepository) *CatalogService {
	return &CatalogService{content: content}
}

func (s *CatalogService) Templates() ([]dto.TaskTemplateItem, error) {
	items, err := s.content.ListTemplates()
	if err != nil {
		return nil, err
	}
	result := make([]dto.TaskTemplateItem, 0, len(items))
	for _, item := range items {
		result = append(result, toTemplateItem(item))
	}
	return result, nil
}

func (s *CatalogService) Template(id uint) (*dto.TaskTemplateItem, error) {
	item, err := s.content.GetTemplate(id)
	if err != nil {
		return nil, support.NewError(http.StatusNotFound, "template_not_found", "task template not found")
	}
	result := toTemplateItem(*item)
	return &result, nil
}

func (s *CatalogService) ApplicationCases() ([]dto.ApplicationCaseItem, error) {
	items, err := s.content.ListApplicationCases()
	if err != nil {
		return nil, err
	}
	result := make([]dto.ApplicationCaseItem, 0, len(items))
	for _, item := range items {
		result = append(result, toApplicationCase(item))
	}
	return result, nil
}

func (s *CatalogService) ApplicationCase(id uint) (*dto.ApplicationCaseItem, error) {
	item, err := s.content.GetApplicationCase(id)
	if err != nil {
		return nil, support.NewError(http.StatusNotFound, "application_case_not_found", "application case not found")
	}
	result := toApplicationCase(*item)
	return &result, nil
}

func (s *CatalogService) DocCategories(docType string) ([]dto.DocumentCategoryItem, error) {
	items, err := s.content.ListDocCategories(docType)
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

func (s *CatalogService) Docs(docType string) ([]dto.DocumentItem, error) {
	items, err := s.content.ListDocuments(docType)
	if err != nil {
		return nil, err
	}
	result := make([]dto.DocumentItem, 0, len(items))
	for _, item := range items {
		result = append(result, toDocumentItem(item))
	}
	return result, nil
}

func (s *CatalogService) Doc(id uint) (*dto.DocumentItem, error) {
	item, err := s.content.GetDocument(id)
	if err != nil {
		return nil, support.NewError(http.StatusNotFound, "doc_not_found", "document not found")
	}
	result := toDocumentItem(*item)
	return &result, nil
}

func (s *CatalogService) FAQs() ([]dto.FAQItem, error) {
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

func (s *CatalogService) Videos() ([]dto.VideoTutorialItem, error) {
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

func (s *CatalogService) DatasetOptions() (*dto.DatasetOptionsResponse, error) {
	templates, err := s.content.ListAgreementTemplates(true)
	if err != nil {
		return nil, err
	}
	privacyOptions, err := s.content.ListDatasetPrivacyOptions(true)
	if err != nil {
		return nil, err
	}
	result := &dto.DatasetOptionsResponse{
		AgreementTemplates: make([]dto.AgreementTemplateItem, 0, len(templates)),
		PrivacyOptions:     make([]dto.DatasetPrivacyOptionItem, 0, len(privacyOptions)),
	}
	for _, item := range templates {
		result.AgreementTemplates = append(result.AgreementTemplates, toAgreementTemplateItem(item))
	}
	for _, item := range privacyOptions {
		result.PrivacyOptions = append(result.PrivacyOptions, toDatasetPrivacyOptionItem(item))
	}
	return result, nil
}
