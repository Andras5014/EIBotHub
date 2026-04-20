package service

import (
	"errors"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/model"
	"github.com/Andras5014/EIBotHub/server/internal/repository"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

type SearchService struct {
	models    *repository.ModelRepository
	datasets  *repository.DatasetRepository
	content   *repository.ContentRepository
	users     *repository.UserRepository
	community *repository.CommunityRepository
}

func NewSearchService(models *repository.ModelRepository, datasets *repository.DatasetRepository, content *repository.ContentRepository, users *repository.UserRepository, community *repository.CommunityRepository) *SearchService {
	return &SearchService{models: models, datasets: datasets, content: content, users: users, community: community}
}

func (s *SearchService) Search(query dto.SearchQuery) (*dto.SearchResponse, error) {
	page, pageSize := normalizePagination(query.Page, query.PageSize)
	results := make([]dto.SearchItem, 0)
	fetchPage := page
	fetchPageSize := pageSize

	if query.Type == "" {
		fetchPage = 1
		fetchPageSize = 50
	}

	if strings.TrimSpace(query.Q) != "" && s.community != nil {
		_ = s.community.AddSearchRecord(&model.SearchRecord{
			Query:      strings.TrimSpace(query.Q),
			SearchType: query.Type,
		})
	}

	if query.Type == "" || query.Type == "model" {
		items, _, err := s.models.Search(query.Q, query.Tags, query.RobotType, fetchPage, fetchPageSize)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			results = append(results, dto.SearchItem{
				ID:        item.ID,
				Type:      "model",
				Title:     item.Name,
				Summary:   item.Summary,
				Tags:      support.SplitCSV(item.Tags),
				Route:     "/models/" + strconv.FormatUint(uint64(item.ID), 10),
				ScoreHint: item.Downloads,
				UpdatedAt: item.UpdatedAt,
			})
		}
	}

	if query.Type == "" || query.Type == "dataset" {
		items, _, err := s.datasets.Search(query.Q, query.Tags, fetchPage, fetchPageSize)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			results = append(results, dto.SearchItem{
				ID:        item.ID,
				Type:      "dataset",
				Title:     item.Name,
				Summary:   item.Summary,
				Tags:      support.SplitCSV(item.Tags),
				Route:     "/datasets/" + strconv.FormatUint(uint64(item.ID), 10),
				ScoreHint: item.Downloads,
				UpdatedAt: item.UpdatedAt,
			})
		}
	}

	if query.Type == "" || query.Type == "task-template" {
		items, _, err := s.content.SearchTemplates(query.Q, fetchPage, fetchPageSize)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			results = append(results, dto.SearchItem{
				ID:        item.ID,
				Type:      "task-template",
				Title:     item.Name,
				Summary:   item.Summary,
				Tags:      []string{item.Category, item.Scene},
				Route:     "/templates/" + strconv.FormatUint(uint64(item.ID), 10),
				ScoreHint: item.UsageCount,
				UpdatedAt: item.UpdatedAt,
			})
		}
	}

	if query.Type == "" || query.Type == "doc" {
		items, _, err := s.content.SearchDocuments(query.Q, fetchPage, fetchPageSize)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			results = append(results, dto.SearchItem{
				ID:        item.ID,
				Type:      "doc",
				Title:     item.Title,
				Summary:   item.Summary,
				Tags:      []string{item.DocType, item.Category.Name},
				Route:     "/docs/" + strconv.FormatUint(uint64(item.ID), 10),
				UpdatedAt: item.UpdatedAt,
			})
		}
	}

	if query.Type == "" || query.Type == "user" {
		userItems, err := s.users.SearchByKeyword(query.Q)
		if err != nil {
			return nil, err
		}
		for _, item := range userItems {
			results = append(results, dto.SearchItem{
				ID:        item.ID,
				Type:      "user",
				Title:     item.Username,
				Summary:   strings.TrimSpace(item.Bio),
				Route:     "/community/users/" + strconv.FormatUint(uint64(item.ID), 10),
				UpdatedAt: item.UpdatedAt,
			})
		}
	}

	if query.Type == "" || query.Type == "skill" {
		items, _, err := s.community.SearchSkills(query.Q, fetchPage, fetchPageSize)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			results = append(results, dto.SearchItem{
				ID:        item.ID,
				Type:      "skill",
				Title:     item.Name,
				Summary:   item.Summary,
				Tags:      []string{item.Category, item.Scene},
				Route:     "/skills/" + strconv.FormatUint(uint64(item.ID), 10),
				ScoreHint: item.UsageCount,
				UpdatedAt: item.UpdatedAt,
			})
		}
	}

	if query.Type == "" || query.Type == "discussion" {
		items, _, err := s.community.SearchDiscussions(query.Q, fetchPage, fetchPageSize)
		if err != nil {
			return nil, err
		}
		for _, item := range items {
			results = append(results, dto.SearchItem{
				ID:        item.ID,
				Type:      "discussion",
				Title:     item.Title,
				Summary:   item.Summary,
				Tags:      []string{item.Category},
				Route:     "/discussions/" + strconv.FormatUint(uint64(item.ID), 10),
				UpdatedAt: item.UpdatedAt,
			})
		}
	}

	results = filterSearchItems(results, query)
	sortSearchItems(results, query.Sort)

	fullResults := append([]dto.SearchItem(nil), results...)
	total := int64(len(fullResults))
	typeCounts := aggregateSearchTypeCounts(fullResults)
	focusType := resolveFocusType(query.Type, fullResults, typeCounts)
	sameTypeItems := pickSearchRecommendations(fullResults, focusType, true, nil, 4)
	relatedItems := pickSearchRecommendations(fullResults, focusType, false, sameTypeItems, 6)
	suggestedQueries, err := s.buildSuggestedQueries(strings.TrimSpace(query.Q), fullResults)
	if err != nil {
		return nil, err
	}
	if query.Type == "" {
		start := (page - 1) * pageSize
		if start >= len(results) {
			return &dto.SearchResponse{
				Items:            []dto.SearchItem{},
				Page:             page,
				PageSize:         pageSize,
				Total:            total,
				FocusType:        focusType,
				TypeCounts:       typeCounts,
				SameTypeItems:    sameTypeItems,
				RelatedItems:     relatedItems,
				SuggestedQueries: suggestedQueries,
			}, nil
		}
		end := start + pageSize
		if end > len(results) {
			end = len(results)
		}
		results = results[start:end]
	}

	return &dto.SearchResponse{
		Items:            results,
		Page:             page,
		PageSize:         pageSize,
		Total:            total,
		FocusType:        focusType,
		TypeCounts:       typeCounts,
		SameTypeItems:    sameTypeItems,
		RelatedItems:     relatedItems,
		SuggestedQueries: suggestedQueries,
	}, nil
}

func (s *SearchService) HotQueries() ([]dto.SearchHotItem, error) {
	curated, err := s.community.ListSearchKeywordConfigs("hot", true)
	if err != nil {
		return nil, err
	}
	organic, err := s.community.HotQueries(8)
	if err != nil {
		return nil, err
	}

	organicCount := make(map[string]int64, len(organic))
	for _, item := range organic {
		organicCount[item.Query] = item.Count
	}

	result := make([]dto.SearchHotItem, 0, 8)
	seen := make(map[string]struct{}, 8)
	for _, item := range curated {
		query := strings.TrimSpace(item.Query)
		if query == "" {
			continue
		}
		if _, exists := seen[query]; exists {
			continue
		}
		result = append(result, dto.SearchHotItem{
			Query: query,
			Count: organicCount[query],
		})
		seen[query] = struct{}{}
		if len(result) >= 8 {
			return result, nil
		}
	}

	for _, item := range organic {
		if _, exists := seen[item.Query]; exists {
			continue
		}
		result = append(result, dto.SearchHotItem{
			Query: item.Query,
			Count: item.Count,
		})
		seen[item.Query] = struct{}{}
		if len(result) >= 8 {
			break
		}
	}
	return result, nil
}

func (s *SearchService) RecommendedQueries() ([]dto.SearchSuggestionItem, error) {
	items, err := s.community.ListSearchKeywordConfigs("recommended", true)
	if err != nil {
		return nil, err
	}
	result := make([]dto.SearchSuggestionItem, 0, len(items))
	for _, item := range items {
		result = append(result, dto.SearchSuggestionItem{Query: item.Query})
	}
	return result, nil
}

func (s *SearchService) AdminSearchKeywords(keywordType string) ([]dto.SearchKeywordConfigItem, error) {
	items, err := s.community.ListSearchKeywordConfigs(keywordType, false)
	if err != nil {
		return nil, err
	}
	result := make([]dto.SearchKeywordConfigItem, 0, len(items))
	for _, item := range items {
		result = append(result, dto.SearchKeywordConfigItem{
			ID:          item.ID,
			Query:       item.Query,
			KeywordType: item.KeywordType,
			SortOrder:   item.SortOrder,
			Enabled:     item.Enabled,
			UpdatedAt:   item.UpdatedAt,
		})
	}
	return result, nil
}

func (s *SearchService) CreateSearchKeyword(input dto.SearchKeywordConfigRequest) (*dto.SearchKeywordConfigItem, error) {
	item := &model.SearchKeywordConfig{
		Query:       strings.TrimSpace(input.Query),
		KeywordType: input.KeywordType,
		SortOrder:   input.SortOrder,
		Enabled:     input.Enabled,
	}
	if err := s.community.CreateSearchKeywordConfig(item); err != nil {
		return nil, err
	}
	return &dto.SearchKeywordConfigItem{
		ID:          item.ID,
		Query:       item.Query,
		KeywordType: item.KeywordType,
		SortOrder:   item.SortOrder,
		Enabled:     item.Enabled,
		UpdatedAt:   item.UpdatedAt,
	}, nil
}

func (s *SearchService) UpdateSearchKeyword(id uint, input dto.SearchKeywordConfigRequest) (*dto.SearchKeywordConfigItem, error) {
	item, err := s.community.GetSearchKeywordConfig(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, support.NewError(http.StatusNotFound, "search_keyword_not_found", "搜索运营词不存在")
		}
		return nil, err
	}
	item.Query = strings.TrimSpace(input.Query)
	item.KeywordType = input.KeywordType
	item.SortOrder = input.SortOrder
	item.Enabled = input.Enabled
	if err := s.community.UpdateSearchKeywordConfig(item); err != nil {
		return nil, err
	}
	return &dto.SearchKeywordConfigItem{
		ID:          item.ID,
		Query:       item.Query,
		KeywordType: item.KeywordType,
		SortOrder:   item.SortOrder,
		Enabled:     item.Enabled,
		UpdatedAt:   item.UpdatedAt,
	}, nil
}

func (s *SearchService) DeleteSearchKeyword(id uint) error {
	item, err := s.community.GetSearchKeywordConfig(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return support.NewError(http.StatusNotFound, "search_keyword_not_found", "搜索运营词不存在")
		}
		return err
	}
	return s.community.DeleteSearchKeywordConfig(item.ID)
}

func filterSearchItems(items []dto.SearchItem, query dto.SearchQuery) []dto.SearchItem {
	filtered := make([]dto.SearchItem, 0, len(items))
	cutoff := time.Time{}
	if query.UpdatedWithin > 0 {
		cutoff = time.Now().AddDate(0, 0, -query.UpdatedWithin)
	}
	requiredTags := splitAndNormalize(query.Tags)

	for _, item := range items {
		if !cutoff.IsZero() && item.UpdatedAt.Before(cutoff) {
			continue
		}
		if len(requiredTags) > 0 && !hasAllTags(item.Tags, requiredTags) {
			continue
		}
		filtered = append(filtered, item)
	}
	return filtered
}

func sortSearchItems(items []dto.SearchItem, sortMode string) {
	sort.SliceStable(items, func(i, j int) bool {
		left := items[i]
		right := items[j]

		switch sortMode {
		case "name":
			return left.Title < right.Title
		case "latest":
			return left.UpdatedAt.After(right.UpdatedAt)
		case "oldest":
			return left.UpdatedAt.Before(right.UpdatedAt)
		case "hot", "downloads":
			if left.ScoreHint == right.ScoreHint {
				return left.UpdatedAt.After(right.UpdatedAt)
			}
			return left.ScoreHint > right.ScoreHint
		default:
			if left.ScoreHint == right.ScoreHint {
				return left.UpdatedAt.After(right.UpdatedAt)
			}
			return left.ScoreHint > right.ScoreHint
		}
	})
}

func aggregateSearchTypeCounts(items []dto.SearchItem) []dto.SearchTypeCountItem {
	counts := make(map[string]int64)
	for _, item := range items {
		counts[item.Type]++
	}
	result := make([]dto.SearchTypeCountItem, 0, len(counts))
	for itemType, count := range counts {
		result = append(result, dto.SearchTypeCountItem{
			Type:  itemType,
			Label: searchTypeLabel(itemType),
			Count: count,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Count == result[j].Count {
			return result[i].Label < result[j].Label
		}
		return result[i].Count > result[j].Count
	})
	return result
}

func resolveFocusType(requestedType string, items []dto.SearchItem, counts []dto.SearchTypeCountItem) string {
	if requestedType != "" {
		return requestedType
	}
	if len(counts) > 0 {
		return counts[0].Type
	}
	if len(items) > 0 {
		return items[0].Type
	}
	return ""
}

func pickSearchRecommendations(items []dto.SearchItem, focusType string, sameType bool, exclude []dto.SearchItem, limit int) []dto.SearchItem {
	if limit <= 0 {
		return []dto.SearchItem{}
	}
	excluded := make(map[string]struct{})
	for _, item := range exclude {
		excluded[searchKey(item)] = struct{}{}
	}
	result := make([]dto.SearchItem, 0, limit)
	for _, item := range items {
		if focusType == "" {
			if sameType {
				continue
			}
		} else {
			if sameType && item.Type != focusType {
				continue
			}
			if !sameType && item.Type == focusType {
				continue
			}
		}
		key := searchKey(item)
		if _, exists := excluded[key]; exists {
			continue
		}
		excluded[key] = struct{}{}
		result = append(result, item)
		if len(result) >= limit {
			break
		}
	}
	return result
}

func (s *SearchService) buildSuggestedQueries(query string, items []dto.SearchItem) ([]dto.SearchSuggestionItem, error) {
	suggestions := make([]dto.SearchSuggestionItem, 0, 6)
	seen := make(map[string]struct{})
	normalizedQuery := strings.ToLower(strings.TrimSpace(query))
	appendSuggestion := func(value string) {
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			return
		}
		lower := strings.ToLower(trimmed)
		if lower == normalizedQuery {
			return
		}
		if _, exists := seen[lower]; exists {
			return
		}
		seen[lower] = struct{}{}
		suggestions = append(suggestions, dto.SearchSuggestionItem{Query: trimmed})
	}

	tagScores := make(map[string]int)
	for _, item := range items {
		for _, tag := range item.Tags {
			trimmed := strings.TrimSpace(tag)
			if trimmed == "" {
				continue
			}
			tagScores[trimmed]++
		}
	}
	type scoredTag struct {
		value string
		score int
	}
	scoredTags := make([]scoredTag, 0, len(tagScores))
	for value, score := range tagScores {
		scoredTags = append(scoredTags, scoredTag{value: value, score: score})
	}
	sort.Slice(scoredTags, func(i, j int) bool {
		if scoredTags[i].score == scoredTags[j].score {
			return scoredTags[i].value < scoredTags[j].value
		}
		return scoredTags[i].score > scoredTags[j].score
	})
	for _, item := range scoredTags {
		appendSuggestion(item.value)
		if len(suggestions) >= 6 {
			return suggestions, nil
		}
	}

	curated, err := s.community.ListSearchKeywordConfigs("recommended", true)
	if err != nil {
		return nil, err
	}
	for _, item := range curated {
		appendSuggestion(item.Query)
		if len(suggestions) >= 6 {
			break
		}
	}
	return suggestions, nil
}

func searchTypeLabel(value string) string {
	return map[string]string{
		"model":         "模型",
		"dataset":       "数据集",
		"task-template": "模板",
		"skill":         "技能",
		"discussion":    "讨论",
		"doc":           "文档",
		"user":          "用户",
	}[value]
}

func searchKey(item dto.SearchItem) string {
	return item.Type + ":" + strconv.FormatUint(uint64(item.ID), 10)
}

func splitAndNormalize(input string) []string {
	parts := strings.Split(input, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.ToLower(strings.TrimSpace(part)); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func hasAllTags(candidate []string, required []string) bool {
	if len(required) == 0 {
		return true
	}
	if len(candidate) == 0 {
		return false
	}
	normalized := make(map[string]struct{}, len(candidate))
	for _, item := range candidate {
		normalized[strings.ToLower(strings.TrimSpace(item))] = struct{}{}
	}
	for _, item := range required {
		if _, ok := normalized[item]; !ok {
			return false
		}
	}
	return true
}
