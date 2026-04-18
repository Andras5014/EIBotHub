package service

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/model"
	"github.com/Andras5014/EIBotHub/server/internal/repository"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

type DatasetCreateInput struct {
	Name          string
	Summary       string
	Description   string
	Tags          []string
	SampleCount   int
	Device        string
	Scene         string
	Privacy       string
	AgreementText string
	Version       string
	Changelog     string
	FilePath      string
	FileName      string
	SamplePreview []string
	Samples       []DatasetSampleInput
}

type DatasetSampleInput struct {
	SampleType  string
	Title       string
	PreviewText string
	FilePath    string
	FileName    string
}

type DatasetService struct {
	repo     *repository.DatasetRepository
	activity *repository.UserActivityRepository
	users    *repository.UserRepository
	storage  *support.LocalStorage
}

func NewDatasetService(repo *repository.DatasetRepository, activity *repository.UserActivityRepository, users *repository.UserRepository, storage *support.LocalStorage) *DatasetService {
	return &DatasetService{repo: repo, activity: activity, users: users, storage: storage}
}

func (s *DatasetService) List(query dto.ResourceListQuery) (map[string]any, error) {
	page, pageSize := normalizePagination(query.Page, query.PageSize)
	items, total, err := s.repo.ListPublished(query.Q, query.Tags, query.Sort, page, pageSize)
	if err != nil {
		return nil, err
	}
	result := make([]dto.ResourceCard, 0, len(items))
	for _, item := range items {
		result = append(result, toDatasetCard(item))
	}
	return support.Paginated(result, page, pageSize, total), nil
}

func (s *DatasetService) Detail(id, userID uint) (*dto.DatasetDetail, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, support.NewError(http.StatusNotFound, "dataset_not_found", "dataset not found")
		}
		return nil, err
	}

	if item.Status != model.StatusPublished && item.OwnerID != userID {
		return nil, support.NewError(http.StatusForbidden, "dataset_unavailable", "dataset is not published")
	}

	favorited := false
	if userID > 0 {
		favorited, err = s.activity.IsFavorited(userID, "dataset", id)
		if err != nil {
			return nil, err
		}
	}

	detail := &dto.DatasetDetail{
		ResourceCard:  toDatasetCard(*item),
		SampleCount:   item.SampleCount,
		Device:        item.Device,
		Scene:         item.Scene,
		Privacy:       item.Privacy,
		AgreementText: item.AgreementText,
		Samples:       make([]dto.DatasetSample, 0, len(item.Samples)),
		Versions:      make([]dto.FileVersion, 0, len(item.Versions)),
		Favorited:     favorited,
	}
	for _, sample := range item.Samples {
		detail.Samples = append(detail.Samples, dto.DatasetSample{
			ID:          sample.ID,
			SampleType:  sample.SampleType,
			Title:       sample.Title,
			PreviewText: sample.PreviewText,
			PreviewURL:  buildStorageURL(sample.FilePath),
			FileName:    sample.FileName,
		})
	}
	for _, version := range item.Versions {
		detail.Versions = append(detail.Versions, dto.FileVersion{
			ID:        version.ID,
			Version:   version.Version,
			FileName:  version.FileName,
			FileURL:   buildStorageURL(version.FilePath),
			Changelog: version.Changelog,
			CreatedAt: version.CreatedAt,
		})
	}
	return detail, nil
}

func (s *DatasetService) Create(ownerID uint, input DatasetCreateInput) (*dto.DatasetDetail, error) {
	resource := &model.Dataset{
		Name:          input.Name,
		Summary:       input.Summary,
		Description:   input.Description,
		Tags:          support.JoinCSV(input.Tags),
		SampleCount:   input.SampleCount,
		Device:        input.Device,
		Scene:         input.Scene,
		Privacy:       input.Privacy,
		AgreementText: input.AgreementText,
		Status:        model.StatusDraft,
		OwnerID:       ownerID,
		Versions: []model.DatasetVersion{
			{
				Version:   firstNonEmpty(input.Version, "v1.0.0"),
				FilePath:  input.FilePath,
				FileName:  input.FileName,
				Changelog: input.Changelog,
			},
		},
		Samples: buildSamples(input.SamplePreview, input.Samples),
	}
	if err := s.repo.Create(resource); err != nil {
		return nil, err
	}
	return s.Detail(resource.ID, ownerID)
}

func (s *DatasetService) AddVersion(id, ownerID uint, input DatasetCreateInput) (*dto.DatasetDetail, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if item.OwnerID != ownerID {
		return nil, support.NewError(http.StatusForbidden, "forbidden", "you do not own this dataset")
	}

	if err := s.repo.CreateVersion(&model.DatasetVersion{
		DatasetID: id,
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

func (s *DatasetService) Submit(id, ownerID uint) error {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if item.OwnerID != ownerID {
		return support.NewError(http.StatusForbidden, "forbidden", "you do not own this dataset")
	}
	item.Status = model.StatusPending
	return s.repo.Update(item)
}

func (s *DatasetService) Download(id, userID uint) error {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if item.Status != model.StatusPublished {
		return support.NewError(http.StatusForbidden, "dataset_unavailable", "dataset is not published")
	}
	if err := s.ensureDatasetAccess(*item, userID); err != nil {
		return err
	}
	accepted, err := s.activity.HasAgreement(userID, id)
	if err != nil {
		return err
	}
	if !accepted {
		return support.NewError(http.StatusBadRequest, "agreement_required", "dataset agreement must be accepted before download")
	}
	if err := s.repo.IncrementDownloads(id); err != nil {
		return err
	}
	return s.activity.AddDownload(&model.DownloadRecord{
		UserID:        userID,
		ResourceType:  "dataset",
		ResourceID:    id,
		ResourceTitle: item.Name,
	})
}

func (s *DatasetService) ConfirmAgreement(id, userID uint) error {
	return s.activity.ConfirmAgreement(&model.AgreementRecord{
		UserID:     userID,
		DatasetID:  id,
		AcceptedAt: time.Now(),
	})
}

func (s *DatasetService) Samples(id, userID uint) ([]dto.DatasetSample, error) {
	item, err := s.Detail(id, userID)
	if err != nil {
		return nil, err
	}
	return item.Samples, nil
}

func (s *DatasetService) CreateDownloadPackageTask(id, userID uint, parts int) (*dto.DownloadPackageTaskItem, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if item.Status != model.StatusPublished {
		return nil, support.NewError(http.StatusForbidden, "dataset_unavailable", "dataset is not published")
	}
	if err := s.ensureDatasetAccess(*item, userID); err != nil {
		return nil, err
	}
	accepted, err := s.activity.HasAgreement(userID, id)
	if err != nil {
		return nil, err
	}
	if !accepted {
		return nil, support.NewError(http.StatusBadRequest, "agreement_required", "dataset agreement must be accepted before download")
	}
	if parts <= 0 {
		parts = 3
	}
	versionFileName := "dataset"
	versionFilePath := ""
	if len(item.Versions) > 0 && item.Versions[0].FileName != "" {
		versionFileName = item.Versions[0].FileName
		versionFilePath = item.Versions[0].FilePath
	}
	if versionFilePath == "" {
		return nil, support.NewError(http.StatusNotFound, "dataset_file_missing", "dataset source file not found")
	}

	partPaths, err := s.storage.SplitFile(versionFilePath, "download-packages", versionFileName, parts)
	if err != nil {
		return nil, err
	}
	partLinks := make([]string, 0, len(partPaths))
	for _, partPath := range partPaths {
		partLinks = append(partLinks, buildStorageURL(partPath))
	}

	bundlePath, err := s.storage.WriteGeneratedFile(
		"download-packages",
		"dataset-"+strconv.FormatUint(uint64(id), 10)+"-bundle.txt",
		"bundle for dataset "+item.Name+"\nsource="+versionFileName+"\nparts="+strconv.Itoa(len(partLinks)),
	)
	if err != nil {
		return nil, err
	}

	task := &model.DownloadPackageTask{
		DatasetID:  id,
		UserID:     userID,
		Status:     "ready",
		BundlePath: bundlePath,
		PartLinks:  support.JoinCSV(partLinks),
		TotalParts: len(partLinks),
	}
	if err := s.repo.CreateDownloadTask(task); err != nil {
		return nil, err
	}
	if err := s.repo.IncrementDownloads(id); err != nil {
		return nil, err
	}
	if err := s.activity.AddDownload(&model.DownloadRecord{
		UserID:        userID,
		ResourceType:  "dataset",
		ResourceID:    id,
		ResourceTitle: item.Name,
	}); err != nil {
		return nil, err
	}
	return &dto.DownloadPackageTaskItem{
		ID:         task.ID,
		DatasetID:  task.DatasetID,
		Status:     task.Status,
		BundleURL:  buildStorageURL(task.BundlePath),
		PartLinks:  partLinks,
		TotalParts: task.TotalParts,
		CreatedAt:  task.CreatedAt,
	}, nil
}

func (s *DatasetService) DownloadPackageTasks(id, userID uint) ([]dto.DownloadPackageTaskItem, error) {
	items, err := s.repo.ListDownloadTasks(id, userID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.DownloadPackageTaskItem, 0, len(items))
	for _, item := range items {
		result = append(result, dto.DownloadPackageTaskItem{
			ID:         item.ID,
			DatasetID:  item.DatasetID,
			Status:     item.Status,
			BundleURL:  buildStorageURL(item.BundlePath),
			PartLinks:  support.SplitCSV(item.PartLinks),
			TotalParts: item.TotalParts,
			CreatedAt:  item.CreatedAt,
		})
	}
	return result, nil
}

func (s *DatasetService) MyAccessRequest(datasetID, userID uint) (*dto.DatasetAccessRequestItem, error) {
	item, err := s.repo.LatestAccessRequest(datasetID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	result := toDatasetAccessRequestItem(*item, "")
	return &result, nil
}

func (s *DatasetService) CreateAccessRequest(datasetID, userID uint, input dto.DatasetAccessRequestPayload) (*dto.DatasetAccessRequestItem, error) {
	dataset, err := s.repo.GetByID(datasetID)
	if err != nil {
		return nil, err
	}
	if dataset.Status != model.StatusPublished {
		return nil, support.NewError(http.StatusForbidden, "dataset_unavailable", "dataset is not published")
	}
	if dataset.Privacy == "public" {
		return nil, support.NewError(http.StatusBadRequest, "dataset_public", "public dataset does not require access request")
	}

	latest, err := s.repo.LatestAccessRequest(datasetID, userID)
	if err == nil {
		if latest.Status == "pending" {
			return nil, support.NewError(http.StatusConflict, "access_request_pending", "access request is already pending")
		}
		if latest.Status == "approved" {
			return nil, support.NewError(http.StatusConflict, "access_request_approved", "access request already approved")
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	item := &model.DatasetAccessRequest{
		DatasetID: datasetID,
		UserID:    userID,
		Reason:    input.Reason,
		Status:    "pending",
	}
	if err := s.repo.CreateAccessRequest(item); err != nil {
		return nil, err
	}
	if err := s.activity.AddNotification(&model.Notification{
		UserID:  dataset.OwnerID,
		Type:    "dataset_access_request",
		Title:   "数据集访问申请",
		Content: "你的数据集收到新的访问申请，请在后台审核。",
	}); err != nil {
		// ignore notification failure
	}
	result := toDatasetAccessRequestItem(*item, "")
	return &result, nil
}

func (s *DatasetService) AdminAccessRequests() ([]dto.DatasetAccessRequestItem, error) {
	items, err := s.repo.ListAccessRequests()
	if err != nil {
		return nil, err
	}
	result := make([]dto.DatasetAccessRequestItem, 0, len(items))
	for _, item := range items {
		userName := ""
		if user, userErr := s.users.FindByID(item.UserID); userErr == nil {
			userName = user.Username
		}
		result = append(result, toDatasetAccessRequestItem(item, userName))
	}
	return result, nil
}

func (s *DatasetService) ReviewAccessRequest(id, reviewerID uint, input dto.DatasetAccessDecisionRequest) error {
	item, err := s.repo.FindAccessRequest(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return support.NewError(http.StatusNotFound, "access_request_not_found", "dataset access request not found")
		}
		return err
	}
	now := time.Now()
	item.Status = input.Decision
	item.ReviewComment = input.Comment
	item.ReviewedBy = &reviewerID
	item.ReviewedAt = &now
	if err := s.repo.UpdateAccessRequest(item); err != nil {
		return err
	}
	content := "你的数据集访问申请已审核。"
	if input.Decision == "approved" {
		content = "你的数据集访问申请已通过。"
	} else {
		content = "你的数据集访问申请已驳回。"
	}
	return s.activity.AddNotification(&model.Notification{
		UserID:  item.UserID,
		Type:    "dataset_access_review",
		Title:   "数据集访问申请结果",
		Content: content,
	})
}

func (s *DatasetService) Mine(ownerID uint) ([]dto.ResourceCard, error) {
	items, err := s.repo.ListByOwner(ownerID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.ResourceCard, 0, len(items))
	for _, item := range items {
		result = append(result, toDatasetCard(item))
	}
	return result, nil
}

func buildSamples(values []string, extra []DatasetSampleInput) []model.DatasetSample {
	result := make([]model.DatasetSample, 0, len(values)+len(extra))
	for index, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			result = append(result, model.DatasetSample{
				SampleType:  "text",
				Title:       "Sample " + strconv.Itoa(index+1),
				PreviewText: trimmed,
			})
		}
	}
	for _, item := range extra {
		if item.Title == "" {
			continue
		}
		result = append(result, model.DatasetSample{
			SampleType:  item.SampleType,
			Title:       item.Title,
			PreviewText: item.PreviewText,
			FilePath:    item.FilePath,
			FileName:    item.FileName,
		})
	}
	return result
}

func (s *DatasetService) ensureDatasetAccess(item model.Dataset, userID uint) error {
	if userID == 0 {
		if item.Privacy == "public" {
			return nil
		}
		return support.NewError(http.StatusForbidden, "dataset_access_login_required", "login required for this dataset")
	}
	if item.OwnerID == userID {
		return nil
	}
	user, err := s.users.FindByID(userID)
	if err == nil && user.Role == model.RoleAdmin {
		return nil
	}
	if item.Privacy == "public" {
		return nil
	}
	request, err := s.repo.LatestAccessRequest(item.ID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return support.NewError(http.StatusForbidden, "dataset_access_request_required", "dataset access approval required")
		}
		return err
	}
	if request.Status != "approved" {
		return support.NewError(http.StatusForbidden, "dataset_access_not_approved", "dataset access has not been approved")
	}
	return nil
}

func toDatasetAccessRequestItem(item model.DatasetAccessRequest, userName string) dto.DatasetAccessRequestItem {
	return dto.DatasetAccessRequestItem{
		ID:            item.ID,
		DatasetID:     item.DatasetID,
		UserID:        item.UserID,
		UserName:      userName,
		Reason:        item.Reason,
		Status:        item.Status,
		ReviewComment: item.ReviewComment,
		ReviewedAt:    item.ReviewedAt,
		CreatedAt:     item.CreatedAt,
	}
}
