package service

import (
	"net/http"
	"slices"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/model"
	"github.com/Andras5014/EIBotHub/server/internal/repository"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

type CatalogService struct {
	models   *repository.ModelRepository
	datasets *repository.DatasetRepository
	content  *repository.ContentRepository
}

func NewCatalogService(models *repository.ModelRepository, datasets *repository.DatasetRepository, content *repository.ContentRepository) *CatalogService {
	return &CatalogService{models: models, datasets: datasets, content: content}
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

func (s *CatalogService) FilterOptions() (*dto.FilterOptionsResponse, error) {
	models, err := s.models.ListAllPublished()
	if err != nil {
		return nil, err
	}
	datasets, err := s.datasets.ListAllPublished()
	if err != nil {
		return nil, err
	}
	templates, err := s.content.ListTemplates()
	if err != nil {
		return nil, err
	}
	cases, err := s.content.ListApplicationCases()
	if err != nil {
		return nil, err
	}
	configs, err := s.content.ListFilterOptionConfigs("", true)
	if err != nil {
		return nil, err
	}

	modelTags := make([]string, 0)
	datasetTags := make([]string, 0)
	robotTypes := make([]string, 0)
	datasetScenes := make([]string, 0)
	templateCategories := make([]string, 0)
	templateScenes := make([]string, 0)
	caseCategories := make([]string, 0)

	for _, item := range models {
		modelTags = append(modelTags, support.SplitCSV(item.Tags)...)
		robotTypes = appendValue(robotTypes, item.RobotType)
	}
	for _, item := range datasets {
		datasetTags = append(datasetTags, support.SplitCSV(item.Tags)...)
		datasetScenes = appendValue(datasetScenes, item.Scene)
	}
	for _, item := range templates {
		templateCategories = appendValue(templateCategories, item.Category)
		templateScenes = appendValue(templateScenes, item.Scene)
	}
	for _, item := range cases {
		caseCategories = appendValue(caseCategories, item.Category)
	}

	configMap := groupFilterConfigs(configs)
	allTags := append(append([]string{}, modelTags...), datasetTags...)
	return &dto.FilterOptionsResponse{
		Tags:                      mergeConfiguredOptions(configMap["tag"], allTags),
		ModelTags:                 mergeConfiguredOptions(configMap["model_tag"], modelTags),
		DatasetTags:               mergeConfiguredOptions(configMap["dataset_tag"], datasetTags),
		RobotTypes:                mergeConfiguredOptions(configMap["robot_type"], robotTypes),
		DatasetScenes:             mergeConfiguredOptions(configMap["dataset_scene"], datasetScenes),
		TemplateCategories:        mergeConfiguredOptions(configMap["template_category"], templateCategories),
		TemplateScenes:            mergeConfiguredOptions(configMap["template_scene"], templateScenes),
		ApplicationCaseCategories: mergeConfiguredOptions(configMap["application_case_category"], caseCategories),
	}, nil
}

func appendValue(values []string, value string) []string {
	if value == "" {
		return values
	}
	return append(values, value)
}

func uniqueSorted(values []string) []string {
	set := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		if value == "" {
			continue
		}
		if _, exists := set[value]; exists {
			continue
		}
		set[value] = struct{}{}
		result = append(result, value)
	}
	slices.Sort(result)
	return result
}

func groupFilterConfigs(items []model.FilterOptionConfig) map[string][]string {
	result := make(map[string][]string)
	for _, item := range items {
		result[item.Kind] = append(result[item.Kind], item.Value)
	}
	return result
}

func mergeConfiguredOptions(configured []string, discovered []string) []string {
	result := make([]string, 0, len(configured)+len(discovered))
	seen := make(map[string]struct{})
	for _, value := range configured {
		if value == "" {
			continue
		}
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	for _, value := range uniqueSorted(discovered) {
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}
