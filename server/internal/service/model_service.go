package service

import (
	"errors"
	"net/http"

	"gorm.io/gorm"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/model"
	"github.com/Andras5014/EIBotHub/server/internal/repository"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

type ModelCreateInput struct {
	Name         string
	Summary      string
	Description  string
	Tags         []string
	RobotType    string
	InputSpec    string
	OutputSpec   string
	License      string
	Dependencies []string
	Version      string
	Changelog    string
	FilePath     string
	FileName     string
}

type ModelService struct {
	repo      *repository.ModelRepository
	activity  *repository.UserActivityRepository
}

func NewModelService(repo *repository.ModelRepository, activity *repository.UserActivityRepository) *ModelService {
	return &ModelService{repo: repo, activity: activity}
}

func (s *ModelService) List(query dto.ResourceListQuery) (map[string]any, error) {
	page, pageSize := normalizePagination(query.Page, query.PageSize)
	items, total, err := s.repo.ListPublished(query.Q, query.Tags, query.RobotType, query.Sort, page, pageSize)
	if err != nil {
		return nil, err
	}

	result := make([]dto.ResourceCard, 0, len(items))
	for _, item := range items {
		result = append(result, toModelCard(item))
	}
	return support.Paginated(result, page, pageSize, total), nil
}

func (s *ModelService) Detail(id, userID uint) (*dto.ModelDetail, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, support.NewError(http.StatusNotFound, "model_not_found", "model not found")
		}
		return nil, err
	}

	if item.Status != model.StatusPublished && item.OwnerID != userID {
		return nil, support.NewError(http.StatusForbidden, "model_unavailable", "model is not published")
	}

	favorited := false
	if userID > 0 {
		favorited, err = s.activity.IsFavorited(userID, "model", id)
		if err != nil {
			return nil, err
		}
	}

	detail := &dto.ModelDetail{
		ResourceCard: toModelCard(*item),
		InputSpec:    item.InputSpec,
		OutputSpec:   item.OutputSpec,
		License:      item.License,
		Dependencies: support.SplitCSV(item.Dependencies),
		Versions:     make([]dto.FileVersion, 0, len(item.Versions)),
		Favorited:    favorited,
	}
	for _, version := range item.Versions {
		detail.Versions = append(detail.Versions, toModelVersion(version))
	}
	return detail, nil
}

func (s *ModelService) Create(ownerID uint, input ModelCreateInput) (*dto.ModelDetail, error) {
	resource := &model.ModelAsset{
		Name:         input.Name,
		Summary:      input.Summary,
		Description:  input.Description,
		Tags:         support.JoinCSV(input.Tags),
		RobotType:    input.RobotType,
		InputSpec:    input.InputSpec,
		OutputSpec:   input.OutputSpec,
		License:      input.License,
		Dependencies: support.JoinCSV(input.Dependencies),
		Status:       model.StatusDraft,
		OwnerID:      ownerID,
		Versions: []model.ModelVersion{
			{
				Version:   firstNonEmpty(input.Version, "v1.0.0"),
				FilePath:  input.FilePath,
				FileName:  input.FileName,
				Changelog: input.Changelog,
			},
		},
	}
	if err := s.repo.Create(resource); err != nil {
		return nil, err
	}
	return s.Detail(resource.ID, ownerID)
}

func (s *ModelService) AddVersion(id, ownerID uint, input ModelCreateInput) (*dto.ModelDetail, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if item.OwnerID != ownerID {
		return nil, support.NewError(http.StatusForbidden, "forbidden", "you do not own this model")
	}

	if err := s.repo.CreateVersion(&model.ModelVersion{
		ModelID:   id,
		Version:   firstNonEmpty(input.Version, "v-next"),
		FilePath:  input.FilePath,
		FileName:  input.FileName,
		Changelog: input.Changelog,
	}); err != nil {
		return nil, err
	}
	item.Status = model.StatusDraft
	if err := s.repo.Update(item); err != nil {
		return nil, err
	}
	return s.Detail(id, ownerID)
}

func (s *ModelService) Submit(id, ownerID uint) error {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if item.OwnerID != ownerID {
		return support.NewError(http.StatusForbidden, "forbidden", "you do not own this model")
	}
	item.Status = model.StatusPending
	return s.repo.Update(item)
}

func (s *ModelService) Download(id, userID uint) error {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if item.Status != model.StatusPublished {
		return support.NewError(http.StatusForbidden, "model_unavailable", "model is not published")
	}

	if err := s.repo.IncrementDownloads(id); err != nil {
		return err
	}
	return s.activity.AddDownload(&model.DownloadRecord{
		UserID:        userID,
		ResourceType:  "model",
		ResourceID:    id,
		ResourceTitle: item.Name,
	})
}

func (s *ModelService) Mine(ownerID uint) ([]dto.ResourceCard, error) {
	items, err := s.repo.ListByOwner(ownerID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.ResourceCard, 0, len(items))
	for _, item := range items {
		result = append(result, toModelCard(item))
	}
	return result, nil
}
