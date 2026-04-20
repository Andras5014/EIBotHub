package service

import (
	"errors"
	"slices"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/model"
	"github.com/Andras5014/EIBotHub/server/internal/repository"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

var homeModuleLabels = map[string]string{
	"hero":          "导读区",
	"models":        "模型推荐",
	"announcements": "公告动态",
	"scenes":        "应用场景",
	"resources":     "数据集与模板",
	"community":     "社区共创",
}

var defaultHomeModuleKeys = []string{"hero", "models", "announcements", "scenes", "resources", "community"}

type PortalService struct {
	portal   *repository.PortalRepository
	models   *repository.ModelRepository
	datasets *repository.DatasetRepository
	content  *repository.ContentRepository
}

func NewPortalService(portal *repository.PortalRepository, models *repository.ModelRepository, datasets *repository.DatasetRepository, content *repository.ContentRepository) *PortalService {
	return &PortalService{portal: portal, models: models, datasets: datasets, content: content}
}

func (s *PortalService) Home() (*dto.HomeResponse, error) {
	moduleSettings, err := s.ModuleSettings()
	if err != nil {
		return nil, err
	}
	heroConfig, err := s.HeroConfig()
	if err != nil {
		return nil, err
	}
	highlights, err := s.HomeHighlights()
	if err != nil {
		return nil, err
	}
	rankingsConfig, err := s.RankingsConfig()
	if err != nil {
		return nil, err
	}
	announcements, err := s.portal.ListAnnouncements()
	if err != nil {
		return nil, err
	}
	featuredResources, err := s.portal.ListFeaturedResources()
	if err != nil {
		return nil, err
	}

	hotModels, err := s.loadFeaturedModels(featuredResources)
	if err != nil {
		return nil, err
	}
	hotDatasets, err := s.loadFeaturedDatasets(featuredResources)
	if err != nil {
		return nil, err
	}
	templates, err := s.loadFeaturedTemplates(featuredResources)
	if err != nil {
		return nil, err
	}
	cases, err := s.loadFeaturedApplicationCases(featuredResources)
	if err != nil {
		return nil, err
	}
	scenePages, err := s.ScenePages()
	if err != nil {
		return nil, err
	}

	result := &dto.HomeResponse{
		PlatformIntro:    heroConfig.Description,
		HeroConfig:       *heroConfig,
		Highlights:       highlights,
		Announcements:    make([]dto.AnnouncementItem, 0, len(announcements)),
		HotModels:        make([]dto.ResourceCard, 0, len(hotModels)),
		HotDatasets:      make([]dto.ResourceCard, 0, len(hotDatasets)),
		TaskTemplates:    make([]dto.TaskTemplateItem, 0, len(templates)),
		ApplicationCases: make([]dto.ApplicationCaseItem, 0, len(cases)),
		ScenePages:       scenePages,
		ModuleSettings:   moduleSettings,
		RankingsConfig:   *rankingsConfig,
	}
	badgeLabels := featuredBadgeLabels(featuredResources)

	for _, item := range announcements {
		result.Announcements = append(result.Announcements, dto.AnnouncementItem{
			ID:          item.ID,
			Title:       item.Title,
			Summary:     item.Summary,
			Link:        item.Link,
			Pinned:      item.Pinned,
			PublishedAt: item.PublishedAt,
		})
	}
	for _, item := range hotModels {
		card := toModelCard(item)
		card.BadgeLabel = badgeLabels["model"][item.ID]
		result.HotModels = append(result.HotModels, card)
	}
	for _, item := range hotDatasets {
		card := toDatasetCard(item)
		card.BadgeLabel = badgeLabels["dataset"][item.ID]
		result.HotDatasets = append(result.HotDatasets, card)
	}
	for _, item := range templates {
		result.TaskTemplates = append(result.TaskTemplates, toTemplateItem(item))
	}
	for _, item := range cases {
		result.ApplicationCases = append(result.ApplicationCases, toApplicationCase(item))
	}
	return result, nil
}

func (s *PortalService) ScenePages() ([]dto.ScenePageItem, error) {
	items, err := s.portal.ListScenePages(true)
	if err != nil {
		return nil, err
	}
	result := make([]dto.ScenePageItem, 0, len(items))
	for _, item := range items {
		result = append(result, toScenePageItem(item))
	}
	return result, nil
}

func (s *PortalService) ScenePageDetail(slug string) (*dto.ScenePageDetail, error) {
	scene, err := s.portal.GetScenePageBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, support.NewError(404, "scene_not_found", "scene page not found")
		}
		return nil, err
	}
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

	result := &dto.ScenePageDetail{
		Scene:            toScenePageItem(*scene),
		Models:           make([]dto.ResourceCard, 0),
		Datasets:         make([]dto.ResourceCard, 0),
		TaskTemplates:    make([]dto.TaskTemplateItem, 0),
		ApplicationCases: make([]dto.ApplicationCaseItem, 0),
	}

	for _, item := range models {
		if matchesScenePage(scene, item.Name, item.Summary, item.Description, item.Tags, item.RobotType) {
			result.Models = append(result.Models, toModelCard(item))
		}
	}
	for _, item := range datasets {
		if matchesScenePage(scene, item.Name, item.Summary, item.Description, item.Tags, item.Scene) {
			result.Datasets = append(result.Datasets, toDatasetCard(item))
		}
	}
	for _, item := range templates {
		if matchesScenePage(scene, item.Name, item.Summary, item.Description, item.Category, item.Scene, item.Guide) {
			result.TaskTemplates = append(result.TaskTemplates, toTemplateItem(item))
		}
	}
	for _, item := range cases {
		if matchesScenePage(scene, item.Title, item.Summary, item.Category, item.Guide) {
			result.ApplicationCases = append(result.ApplicationCases, toApplicationCase(item))
		}
	}
	return result, nil
}

func (s *PortalService) ModuleSettings() ([]dto.ModuleSettingItem, error) {
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
			sortOrder := item.SortOrder
			if sortOrder == 0 {
				sortOrder = defaultModuleSortOrders()[key]
			}
			result = append(result, dto.ModuleSettingItem{
				ID:        item.ID,
				ModuleKey: key,
				Label:     homeModuleLabels[key],
				SortOrder: sortOrder,
				Enabled:   item.Enabled,
				UpdatedAt: item.UpdatedAt,
			})
			continue
		}
		result = append(result, dto.ModuleSettingItem{
			ModuleKey: key,
			Label:     homeModuleLabels[key],
			SortOrder: defaultModuleSortOrders()[key],
			Enabled:   true,
		})
	}
	return result, nil
}

func (s *PortalService) HomeHighlights() ([]string, error) {
	items, err := s.portal.ListHomeHighlights(true)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return []string{
			"统一管理模型、数据集和任务模板",
			"面向移动作业、巡逻巡检、搬运等场景",
			"提供文档、下载统计、审核和运营后台",
		}, nil
	}
	result := make([]string, 0, len(items))
	for _, item := range items {
		result = append(result, item.Text)
	}
	return result, nil
}

func (s *PortalService) HeroConfig() (*dto.HomeHeroConfigItem, error) {
	item, err := s.portal.GetHomeHeroConfig()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || item == nil {
			return defaultHeroConfig(), nil
		}
		return nil, err
	}
	return &dto.HomeHeroConfigItem{
		ID:              item.ID,
		Tagline:         item.Tagline,
		Title:           item.Title,
		Description:     item.Description,
		PrimaryButton:   item.PrimaryButton,
		SecondaryButton: item.SecondaryButton,
		SearchButton:    item.SearchButton,
		UpdatedAt:       item.UpdatedAt,
	}, nil
}

func (s *PortalService) FeaturedResources() ([]dto.FeaturedResourceItem, error) {
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
			BadgeLabel:   item.BadgeLabel,
			SortOrder:    item.SortOrder,
			Enabled:      item.Enabled,
			UpdatedAt:    item.UpdatedAt,
		})
	}
	return result, nil
}

func (s *PortalService) RankingsConfig() (*dto.RankingConfigItem, error) {
	item, err := s.portal.GetRankingConfig()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || item == nil {
			return defaultRankingConfig(), nil
		}
		return nil, err
	}
	return &dto.RankingConfigItem{
		ID:        item.ID,
		Title:     item.Title,
		Subtitle:  item.Subtitle,
		Limit:     item.Limit,
		Enabled:   item.Enabled,
		UpdatedAt: item.UpdatedAt,
	}, nil
}

func (s *PortalService) resolveFeaturedResource(item model.FeaturedResource) (string, string, string) {
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

func (s *PortalService) loadFeaturedModels(featured []model.FeaturedResource) ([]model.ModelAsset, error) {
	ids := featuredIDs(featured, "model")
	if len(ids) == 0 {
		return s.models.Top(3)
	}
	items, err := s.models.ListPublishedByIDs(ids)
	if err != nil {
		return nil, err
	}
	return orderByIDs(items, ids, func(item model.ModelAsset) uint { return item.ID }), nil
}

func (s *PortalService) loadFeaturedDatasets(featured []model.FeaturedResource) ([]model.Dataset, error) {
	ids := featuredIDs(featured, "dataset")
	if len(ids) == 0 {
		return s.datasets.Top(3)
	}
	items, err := s.datasets.ListPublishedByIDs(ids)
	if err != nil {
		return nil, err
	}
	return orderByIDs(items, ids, func(item model.Dataset) uint { return item.ID }), nil
}

func (s *PortalService) loadFeaturedTemplates(featured []model.FeaturedResource) ([]model.TaskTemplate, error) {
	ids := featuredIDs(featured, "task-template")
	if len(ids) == 0 {
		items, err := s.content.ListTemplates()
		if err != nil {
			return nil, err
		}
		if len(items) > 5 {
			items = items[:5]
		}
		return items, nil
	}
	items, err := s.content.ListTemplatesByIDs(ids)
	if err != nil {
		return nil, err
	}
	return orderByIDs(items, ids, func(item model.TaskTemplate) uint { return item.ID }), nil
}

func (s *PortalService) loadFeaturedApplicationCases(featured []model.FeaturedResource) ([]model.ApplicationCase, error) {
	ids := featuredIDs(featured, "application-case")
	if len(ids) == 0 {
		items, err := s.content.ListApplicationCases()
		if err != nil {
			return nil, err
		}
		if len(items) > 5 {
			items = items[:5]
		}
		return items, nil
	}
	items, err := s.content.ListApplicationCasesByIDs(ids)
	if err != nil {
		return nil, err
	}
	return orderByIDs(items, ids, func(item model.ApplicationCase) uint { return item.ID }), nil
}

func featuredIDs(items []model.FeaturedResource, resourceType string) []uint {
	ids := make([]uint, 0)
	seen := make(map[uint]struct{})
	for _, item := range items {
		if item.ResourceType == resourceType && item.Enabled {
			if _, exists := seen[item.ResourceID]; exists {
				continue
			}
			ids = append(ids, item.ResourceID)
			seen[item.ResourceID] = struct{}{}
		}
	}
	return ids
}

func featuredBadgeLabels(items []model.FeaturedResource) map[string]map[uint]string {
	result := map[string]map[uint]string{
		"model":   {},
		"dataset": {},
	}
	for _, item := range items {
		if !item.Enabled || item.BadgeLabel == "" {
			continue
		}
		if _, ok := result[item.ResourceType]; !ok {
			continue
		}
		if _, exists := result[item.ResourceType][item.ResourceID]; exists {
			continue
		}
		result[item.ResourceType][item.ResourceID] = item.BadgeLabel
	}
	return result
}

func orderByIDs[T any](items []T, ids []uint, getID func(T) uint) []T {
	if len(ids) == 0 || len(items) == 0 {
		return items
	}
	ordered := make([]T, 0, len(items))
	for _, id := range ids {
		index := slices.IndexFunc(items, func(item T) bool {
			return getID(item) == id
		})
		if index >= 0 {
			ordered = append(ordered, items[index])
		}
	}
	return ordered
}

func uintToString(value uint) string {
	return strconv.FormatUint(uint64(value), 10)
}

func defaultModuleSortOrders() map[string]int {
	return map[string]int{
		"hero":          10,
		"models":        20,
		"announcements": 30,
		"scenes":        40,
		"resources":     50,
		"community":     60,
	}
}

func defaultRankingConfig() *dto.RankingConfigItem {
	return &dto.RankingConfigItem{
		Title:    "贡献排行榜",
		Subtitle: "基于积分展示近期社区贡献活跃度。",
		Limit:    5,
		Enabled:  true,
	}
}

func defaultHeroConfig() *dto.HomeHeroConfigItem {
	return &dto.HomeHeroConfigItem{
		Tagline:         "OpenLoong 风格信息门户",
		Title:           "围绕模型、数据集、模板与文档的开放社区",
		Description:     "开放社区聚合模型、数据集、任务模板、文档与具身应用案例，帮助开发者围绕机器人场景快速完成接入、训练、部署和复用。",
		PrimaryButton:   "上传模型",
		SecondaryButton: "上传数据集",
		SearchButton:    "进入全局搜索",
	}
}

func toScenePageItem(item model.ScenePageConfig) dto.ScenePageItem {
	return dto.ScenePageItem{
		ID:          item.ID,
		Slug:        item.Slug,
		Name:        item.Name,
		Tagline:     item.Tagline,
		Summary:     item.Summary,
		Description: item.Description,
		SortOrder:   item.SortOrder,
		Enabled:     item.Enabled,
		UpdatedAt:   item.UpdatedAt,
	}
}

func matchesScenePage(scene *model.ScenePageConfig, values ...string) bool {
	keywords := []string{scene.Name, scene.Tagline, scene.Summary, scene.Description, scene.Slug}
	for _, value := range values {
		lowerValue := strings.ToLower(value)
		for _, keyword := range keywords {
			trimmed := strings.TrimSpace(keyword)
			if trimmed == "" {
				continue
			}
			if strings.Contains(lowerValue, strings.ToLower(trimmed)) {
				return true
			}
		}
	}
	return false
}
