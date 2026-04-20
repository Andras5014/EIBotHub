package service

import (
	"errors"
	"net/http"
	"path/filepath"
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

type DatasetUpdateInput struct {
	Name          string
	Summary       string
	Description   string
	Tags          []string
	SampleCount   int
	Device        string
	Scene         string
	Privacy       string
	AgreementText string
	SamplePreview []string
}

type DatasetSampleInput struct {
	SampleType  string
	Title       string
	PreviewText string
	FilePath    string
	FileName    string
}

type DatasetService struct {
	repo          *repository.DatasetRepository
	activity      *repository.UserActivityRepository
	users         *repository.UserRepository
	verifications *repository.VerificationRepository
	reviews       *repository.ReviewRepository
	storage       support.ObjectStorage
	files         *FileService
}

func NewDatasetService(repo *repository.DatasetRepository, activity *repository.UserActivityRepository, users *repository.UserRepository, verifications *repository.VerificationRepository, reviews *repository.ReviewRepository, storage support.ObjectStorage, files *FileService) *DatasetService {
	return &DatasetService{repo: repo, activity: activity, users: users, verifications: verifications, reviews: reviews, storage: storage, files: files}
}

func (s *DatasetService) List(query dto.ResourceListQuery) (map[string]any, error) {
	page, pageSize := normalizePagination(query.Page, query.PageSize)
	items, total, err := s.repo.ListPublished(query.Q, query.Tags, query.Scene, query.Sort, page, pageSize)
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
			PreviewURL:  buildAuthorizedFileURL(s.files, sample.FilePath, sample.FileName, item.OwnerID, true),
			FileName:    sample.FileName,
		})
	}
	for _, version := range item.Versions {
		detail.Versions = append(detail.Versions, dto.FileVersion{
			ID:        version.ID,
			Version:   version.Version,
			FileName:  version.FileName,
			FileURL:   buildAuthorizedFileURL(s.files, version.FilePath, version.FileName, item.OwnerID, false),
			Changelog: version.Changelog,
			CreatedAt: version.CreatedAt,
		})
	}
	if reviewComment, ok, reviewErr := latestReviewComment(s.reviews, "datasets", id, item.Status, item.OwnerID == userID); reviewErr != nil {
		return nil, reviewErr
	} else if ok {
		detail.ReviewComment = reviewComment
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

func (s *DatasetService) Update(id, ownerID uint, input DatasetUpdateInput) (*dto.DatasetDetail, error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, support.NewError(http.StatusNotFound, "dataset_not_found", "dataset not found")
		}
		return nil, err
	}
	if item.OwnerID != ownerID {
		return nil, support.NewError(http.StatusForbidden, "forbidden", "you do not own this dataset")
	}

	item.Name = input.Name
	item.Summary = input.Summary
	item.Description = input.Description
	item.Tags = support.JoinCSV(input.Tags)
	item.SampleCount = input.SampleCount
	item.Device = input.Device
	item.Scene = input.Scene
	item.Privacy = input.Privacy
	item.AgreementText = input.AgreementText
	item.Status = model.StatusDraft

	if err := s.repo.Update(item); err != nil {
		return nil, err
	}
	if err := s.repo.ReplaceTextSamples(id, buildTextSamples(input.SamplePreview)); err != nil {
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
	accessRequest, err := s.ensureDatasetAccess(*item, userID)
	if err != nil {
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
	if err := s.consumeDatasetAccess(accessRequest); err != nil {
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
	accessRequest, err := s.ensureDatasetAccess(*item, userID)
	if err != nil {
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
		partLinks = append(partLinks, buildAuthorizedFileURL(s.files, partPath, filepath.Base(partPath), item.OwnerID, false))
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
	if err := s.consumeDatasetAccess(accessRequest); err != nil {
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
		BundleURL:  buildAuthorizedFileURL(s.files, task.BundlePath, filepath.Base(task.BundlePath), item.OwnerID, false),
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
			BundleURL:  buildAuthorizedFileURL(s.files, item.BundlePath, filepath.Base(item.BundlePath), 0, false),
			PartLinks:  s.authorizePartLinks(support.SplitCSV(item.PartLinks)),
			TotalParts: item.TotalParts,
			CreatedAt:  item.CreatedAt,
		})
	}
	return result, nil
}

func (s *DatasetService) authorizePartLinks(partPaths []string) []string {
	if len(partPaths) == 0 {
		return []string{}
	}
	result := make([]string, 0, len(partPaths))
	for _, partPath := range partPaths {
		result = append(result, buildAuthorizedFileURL(s.files, partPath, filepath.Base(partPath), 0, false))
	}
	return result
}

func (s *DatasetService) MyAccessRequest(datasetID, userID uint) (*dto.DatasetAccessRequestItem, error) {
	item, err := s.repo.LatestAccessRequest(datasetID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	result, err := s.datasetAccessRequestItems([]model.DatasetAccessRequest{*item}, false)
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, nil
	}
	return &result[0], nil
}

func (s *DatasetService) MyAccessRequests(userID uint) ([]dto.DatasetAccessRequestItem, error) {
	items, err := s.repo.ListAccessRequestsByUser(userID)
	if err != nil {
		return nil, err
	}
	return s.datasetAccessRequestItems(items, false)
}

func (s *DatasetService) MyAccessRequestHistory(datasetID, userID uint) ([]dto.DatasetAccessRequestItem, error) {
	items, err := s.repo.ListAccessRequestsByDatasetAndUser(datasetID, userID)
	if err != nil {
		return nil, err
	}
	return s.datasetAccessRequestItems(items, false)
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
	if dataset.Privacy == "internal" {
		approved, err := s.hasApprovedVerification(userID)
		if err != nil {
			return nil, err
		}
		if approved {
			return nil, support.NewError(http.StatusConflict, "dataset_internal_auto_access", "approved developer verification already grants internal dataset access")
		}
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
		DatasetID:         datasetID,
		UserID:            userID,
		Reason:            input.Reason,
		Status:            "pending",
		RequiredApprovals: datasetRequiredApprovals(dataset.Privacy),
	}
	if err := s.repo.CreateAccessRequest(item); err != nil {
		return nil, err
	}
	s.notifyDatasetAccessRequestCreated(*dataset, *item)
	result, err := s.datasetAccessRequestItems([]model.DatasetAccessRequest{*item}, false)
	if err != nil {
		return nil, err
	}
	return &result[0], nil
}

func (s *DatasetService) AdminAccessRequests(query dto.DatasetAccessAdminQuery) ([]dto.DatasetAccessRequestItem, error) {
	items, err := s.repo.ListAccessRequests()
	if err != nil {
		return nil, err
	}
	result, err := s.datasetAccessRequestItems(items, true)
	if err != nil {
		return nil, err
	}
	return filterDatasetAccessRequestItems(result, query), nil
}

func (s *DatasetService) BatchReviewAccessRequests(reviewerID uint, input dto.BatchDatasetAccessDecisionRequest) error {
	items, err := s.repo.FindAccessRequestsByIDs(input.IDs)
	if err != nil {
		return err
	}

	uniqueIDs := make(map[uint]struct{}, len(input.IDs))
	for _, id := range input.IDs {
		uniqueIDs[id] = struct{}{}
	}
	if len(items) != len(uniqueIDs) {
		return support.NewError(http.StatusNotFound, "access_request_not_found", "dataset access request not found")
	}

	now := time.Now()
	for index := range items {
		if items[index].Status != "pending" {
			return support.NewError(http.StatusConflict, "access_request_not_pending", "only pending requests can be reviewed in batch")
		}
		if input.Decision == "approved" {
			if err := progressDatasetAccessApproval(&items[index], reviewerID, input.Comment, now, input.ValidDays, input.DownloadLimit); err != nil {
				return err
			}
			continue
		}
		rejectDatasetAccessApproval(&items[index], reviewerID, input.Comment, now)
	}

	if err := s.repo.UpdateAccessRequests(items); err != nil {
		return err
	}

	for _, item := range items {
		if err := s.activity.AddNotification(&model.Notification{
			UserID:  item.UserID,
			Type:    "dataset_access_review",
			Title:   "数据集访问申请结果",
			Content: datasetAccessReviewNotificationContent(item),
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s *DatasetService) datasetAccessRequestItems(items []model.DatasetAccessRequest, includeUserName bool) ([]dto.DatasetAccessRequestItem, error) {
	datasetIDs := make([]uint, 0, len(items))
	seenDatasetIDs := make(map[uint]struct{}, len(items))
	for _, item := range items {
		if _, exists := seenDatasetIDs[item.DatasetID]; exists {
			continue
		}
		seenDatasetIDs[item.DatasetID] = struct{}{}
		datasetIDs = append(datasetIDs, item.DatasetID)
	}

	datasets, err := s.repo.ListByIDs(datasetIDs)
	if err != nil {
		return nil, err
	}
	datasetMeta := make(map[uint]dto.DatasetAccessRequestItem, len(datasets))
	for _, item := range datasets {
		datasetMeta[item.ID] = dto.DatasetAccessRequestItem{
			DatasetName:      item.Name,
			DatasetPrivacy:   item.Privacy,
			DatasetOwnerID:   item.OwnerID,
			DatasetOwnerName: item.Owner.Username,
		}
	}

	result := make([]dto.DatasetAccessRequestItem, 0, len(items))
	for _, item := range items {
		userName := ""
		if includeUserName {
			if user, userErr := s.users.FindByID(item.UserID); userErr == nil {
				userName = user.Username
			}
		}
		result = append(result, toDatasetAccessRequestItem(item, userName, datasetMeta[item.DatasetID]))
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
	if item.Status != "pending" {
		return support.NewError(http.StatusConflict, "access_request_not_pending", "only pending requests can be reviewed")
	}
	now := time.Now()
	if input.Decision == "approved" {
		if err := progressDatasetAccessApproval(item, reviewerID, input.Comment, now, input.ValidDays, input.DownloadLimit); err != nil {
			return err
		}
	} else {
		rejectDatasetAccessApproval(item, reviewerID, input.Comment, now)
	}
	if err := s.repo.UpdateAccessRequest(item); err != nil {
		return err
	}
	return s.activity.AddNotification(&model.Notification{
		UserID:  item.UserID,
		Type:    "dataset_access_review",
		Title:   "数据集访问申请结果",
		Content: datasetAccessReviewNotificationContent(*item),
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
	result := buildTextSamples(values)
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

func buildTextSamples(values []string) []model.DatasetSample {
	result := make([]model.DatasetSample, 0, len(values))
	for index, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			result = append(result, model.DatasetSample{
				SampleType:  "text",
				Title:       "Sample " + strconv.Itoa(index+1),
				PreviewText: trimmed,
			})
		}
	}
	return result
}

func (s *DatasetService) ensureDatasetAccess(item model.Dataset, userID uint) (*model.DatasetAccessRequest, error) {
	if userID == 0 {
		if item.Privacy == "public" {
			return nil, nil
		}
		return nil, support.NewError(http.StatusForbidden, "dataset_access_login_required", "login required for this dataset")
	}
	if item.OwnerID == userID {
		return nil, nil
	}
	user, err := s.users.FindByID(userID)
	if err == nil && model.IsAdminRole(user.Role) {
		return nil, nil
	}
	if item.Privacy == "public" {
		return nil, nil
	}
	if item.Privacy == "internal" {
		approved, verificationErr := s.hasApprovedVerification(userID)
		if verificationErr != nil {
			return nil, verificationErr
		}
		if approved {
			return nil, nil
		}
	}
	request, err := s.repo.LatestAccessRequest(item.ID, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, support.NewError(http.StatusForbidden, "dataset_access_request_required", "dataset access approval required")
		}
		return nil, err
	}
	if request.Status != "approved" {
		return nil, support.NewError(http.StatusForbidden, "dataset_access_not_approved", "dataset access has not been approved")
	}
	if request.ApprovalExpiresAt != nil && time.Now().After(*request.ApprovalExpiresAt) {
		return nil, support.NewError(http.StatusForbidden, "dataset_access_expired", "dataset access approval has expired")
	}
	if request.DownloadLimit > 0 && request.DownloadCount >= request.DownloadLimit {
		return nil, support.NewError(http.StatusForbidden, "dataset_access_limit_exceeded", "dataset access download limit reached")
	}
	return request, nil
}

func (s *DatasetService) hasApprovedVerification(userID uint) (bool, error) {
	if s.verifications == nil {
		return false, nil
	}
	item, err := s.verifications.LatestByUser(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return item.Status == "approved", nil
}

func (s *DatasetService) consumeDatasetAccess(request *model.DatasetAccessRequest) error {
	if request == nil || request.DownloadLimit == 0 {
		return nil
	}
	request.DownloadCount += 1
	return s.repo.UpdateAccessRequest(request)
}

func applyDatasetAccessGrant(item *model.DatasetAccessRequest, decision string, now time.Time, validDays, downloadLimit int) {
	if decision != "approved" {
		item.ApprovalExpiresAt = nil
		item.DownloadLimit = 0
		item.DownloadCount = 0
		return
	}
	if validDays > 0 {
		expiresAt := now.AddDate(0, 0, validDays)
		item.ApprovalExpiresAt = &expiresAt
	} else {
		item.ApprovalExpiresAt = nil
	}
	item.DownloadLimit = downloadLimit
	item.DownloadCount = 0
}

func progressDatasetAccessApproval(item *model.DatasetAccessRequest, reviewerID uint, comment string, now time.Time, validDays, downloadLimit int) error {
	requiredApprovals := item.RequiredApprovals
	if requiredApprovals <= 0 {
		requiredApprovals = 1
		item.RequiredApprovals = requiredApprovals
	}
	approverIDs := parseUintCSV(item.ApproverIDs)
	for _, approverID := range approverIDs {
		if approverID == reviewerID {
			return support.NewError(http.StatusConflict, "access_request_stage_already_approved", "current reviewer has already approved this request")
		}
	}
	approverIDs = append(approverIDs, reviewerID)
	item.ApproverIDs = joinUintCSV(approverIDs)
	item.ApprovalStage += 1
	item.ReviewComment = comment
	item.ReviewedBy = &reviewerID
	item.ReviewedAt = &now
	if item.ApprovalStage >= requiredApprovals {
		item.Status = "approved"
		applyDatasetAccessGrant(item, "approved", now, validDays, downloadLimit)
		return nil
	}
	item.Status = "pending"
	item.ApprovalExpiresAt = nil
	item.DownloadLimit = 0
	item.DownloadCount = 0
	return nil
}

func rejectDatasetAccessApproval(item *model.DatasetAccessRequest, reviewerID uint, comment string, now time.Time) {
	item.Status = "rejected"
	item.ReviewComment = comment
	item.ReviewedBy = &reviewerID
	item.ReviewedAt = &now
	applyDatasetAccessGrant(item, "rejected", now, 0, 0)
}

func datasetAccessReviewNotificationContent(item model.DatasetAccessRequest) string {
	if item.Status == "rejected" {
		return "你的数据集访问申请已驳回。"
	}
	if item.Status == "approved" {
		return "你的数据集访问申请已通过。"
	}
	return "你的数据集访问申请已完成第 " + strconv.Itoa(item.ApprovalStage) + "/" + strconv.Itoa(item.RequiredApprovals) + " 级审批，仍需继续审核。"
}

func datasetRequiredApprovals(privacy string) int {
	if privacy == "restricted" {
		return 2
	}
	return 1
}

func filterDatasetAccessRequestItems(items []dto.DatasetAccessRequestItem, query dto.DatasetAccessAdminQuery) []dto.DatasetAccessRequestItem {
	keyword := strings.ToLower(strings.TrimSpace(query.Q))
	result := make([]dto.DatasetAccessRequestItem, 0, len(items))
	for _, item := range items {
		if query.Status != "" && item.Status != query.Status {
			continue
		}
		if query.Privacy != "" && item.DatasetPrivacy != query.Privacy {
			continue
		}
		if query.OwnerID > 0 && item.DatasetOwnerID != query.OwnerID {
			continue
		}
		if query.SLAStatus == "overdue" && (!item.SLAOverdue || item.Status != "pending") {
			continue
		}
		if query.SLAStatus == "ontrack" && (item.Status != "pending" || item.SLAOverdue) {
			continue
		}
		if keyword != "" && !matchesDatasetAccessKeyword(keyword, item) {
			continue
		}
		result = append(result, item)
	}
	return result
}

func matchesDatasetAccessKeyword(keyword string, item dto.DatasetAccessRequestItem) bool {
	if keyword == "" {
		return true
	}
	content := strings.ToLower(strings.Join([]string{
		item.DatasetName,
		item.DatasetPrivacy,
		item.DatasetOwnerName,
		item.UserName,
		item.Reason,
		item.ReviewComment,
	}, " "))
	return strings.Contains(content, keyword)
}

func toDatasetAccessRequestItem(item model.DatasetAccessRequest, userName string, meta dto.DatasetAccessRequestItem) dto.DatasetAccessRequestItem {
	requiredApprovals := item.RequiredApprovals
	if requiredApprovals <= 0 {
		requiredApprovals = datasetRequiredApprovals(meta.DatasetPrivacy)
	}
	approvalStage := item.ApprovalStage
	if item.Status == "approved" && approvalStage < requiredApprovals {
		approvalStage = requiredApprovals
	}
	isExpired := item.ApprovalExpiresAt != nil && time.Now().After(*item.ApprovalExpiresAt)
	slaHours := datasetAccessSLAHours(meta.DatasetPrivacy)
	var slaDeadlineAt *time.Time
	slaOverdue := false
	slaRemainingMinutes := 0
	if slaHours > 0 {
		deadline := item.CreatedAt.Add(time.Duration(slaHours) * time.Hour)
		slaDeadlineAt = &deadline
		if item.Status == "pending" {
			slaRemainingMinutes = int(time.Until(deadline).Minutes())
			slaOverdue = slaRemainingMinutes < 0
		}
	}
	var remainingDownloads *int
	if item.DownloadLimit > 0 {
		remaining := item.DownloadLimit - item.DownloadCount
		if remaining < 0 {
			remaining = 0
		}
		remainingDownloads = &remaining
	}
	authorizationActive := item.Status == "approved" && !isExpired && (item.DownloadLimit == 0 || item.DownloadCount < item.DownloadLimit)
	return dto.DatasetAccessRequestItem{
		ID:                  item.ID,
		DatasetID:           item.DatasetID,
		DatasetName:         meta.DatasetName,
		DatasetPrivacy:      meta.DatasetPrivacy,
		DatasetOwnerID:      meta.DatasetOwnerID,
		DatasetOwnerName:    meta.DatasetOwnerName,
		UserID:              item.UserID,
		UserName:            userName,
		Reason:              item.Reason,
		Status:              item.Status,
		ReviewComment:       item.ReviewComment,
		ApprovalStage:       approvalStage,
		RequiredApprovals:   requiredApprovals,
		ApprovalExpiresAt:   item.ApprovalExpiresAt,
		DownloadLimit:       item.DownloadLimit,
		DownloadCount:       item.DownloadCount,
		RemainingDownloads:  remainingDownloads,
		IsExpired:           isExpired,
		AuthorizationActive: authorizationActive,
		SLAHours:            slaHours,
		SLADeadlineAt:       slaDeadlineAt,
		SLAOverdue:          slaOverdue,
		SLARemainingMinutes: slaRemainingMinutes,
		ReviewedAt:          item.ReviewedAt,
		CreatedAt:           item.CreatedAt,
	}
}

func parseUintCSV(input string) []uint {
	if strings.TrimSpace(input) == "" {
		return nil
	}
	parts := strings.Split(input, ",")
	result := make([]uint, 0, len(parts))
	for _, part := range parts {
		value, err := strconv.ParseUint(strings.TrimSpace(part), 10, 64)
		if err != nil {
			continue
		}
		result = append(result, uint(value))
	}
	return result
}

func joinUintCSV(values []uint) string {
	if len(values) == 0 {
		return ""
	}
	parts := make([]string, 0, len(values))
	for _, value := range values {
		parts = append(parts, strconv.FormatUint(uint64(value), 10))
	}
	return strings.Join(parts, ",")
}

func datasetAccessSLAHours(privacy string) int {
	switch privacy {
	case "internal":
		return 24
	case "restricted":
		return 48
	default:
		return 0
	}
}

func (s *DatasetService) notifyDatasetAccessRequestCreated(dataset model.Dataset, request model.DatasetAccessRequest) {
	if s.activity == nil {
		return
	}
	slaHours := datasetAccessSLAHours(dataset.Privacy)
	slaText := ""
	if slaHours > 0 {
		slaText = "，SLA " + strconv.Itoa(slaHours) + " 小时"
	}
	_ = s.activity.AddNotification(&model.Notification{
		UserID:  request.UserID,
		Type:    "dataset_access_request_submitted",
		Title:   "访问申请已提交",
		Content: "你已提交对数据集《" + dataset.Name + "》的访问申请" + slaText + "。",
	})

	recipients := map[uint]struct{}{
		dataset.OwnerID: {},
	}
	if s.users != nil {
		if admins, err := s.users.ListByRoles(model.RoleAdmin, model.RoleSuperAdmin); err == nil {
			for _, admin := range admins {
				recipients[admin.ID] = struct{}{}
			}
		}
	}
	for recipientID := range recipients {
		_ = s.activity.AddNotification(&model.Notification{
			UserID:  recipientID,
			Type:    "dataset_access_request",
			Title:   "数据集访问申请",
			Content: "数据集《" + dataset.Name + "》收到新的访问申请" + slaText + "，请尽快审核。",
		})
	}
}
