package service

import (
	"sort"
	"strconv"
	"strings"
	"time"

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

func (s *SearchService) Search(query dto.SearchQuery) (map[string]any, error) {
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

	total := int64(len(results))
	if query.Type == "" {
		start := (page - 1) * pageSize
		if start >= len(results) {
			return support.Paginated([]dto.SearchItem{}, page, pageSize, total), nil
		}
		end := start + pageSize
		if end > len(results) {
			end = len(results)
		}
		results = results[start:end]
	}

	return support.Paginated(results, page, pageSize, total), nil
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
