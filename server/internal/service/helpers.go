package service

import (
	"path"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/model"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

func normalizePagination(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 12
	}
	return page, pageSize
}

func firstNonEmpty(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func buildStorageURL(filePath string) string {
	if filePath == "" {
		return ""
	}
	return "/" + path.Clean("storage/"+filePath)
}

func toModelCard(item model.ModelAsset) dto.ResourceCard {
	return dto.ResourceCard{
		ID:          item.ID,
		Name:        item.Name,
		Summary:     item.Summary,
		Description: item.Description,
		Type:        "model",
		Tags:        support.SplitCSV(item.Tags),
		RobotType:   item.RobotType,
		Downloads:   item.Downloads,
		Status:      item.Status,
		Owner:       item.Owner.Username,
		UpdatedAt:   item.UpdatedAt,
	}
}

func toDatasetCard(item model.Dataset) dto.ResourceCard {
	return dto.ResourceCard{
		ID:          item.ID,
		Name:        item.Name,
		Summary:     item.Summary,
		Description: item.Description,
		Type:        "dataset",
		Tags:        support.SplitCSV(item.Tags),
		Downloads:   item.Downloads,
		Status:      item.Status,
		Owner:       item.Owner.Username,
		UpdatedAt:   item.UpdatedAt,
	}
}

func toModelVersion(version model.ModelVersion) dto.FileVersion {
	return dto.FileVersion{
		ID:        version.ID,
		Version:   version.Version,
		FileName:  version.FileName,
		FileURL:   buildStorageURL(version.FilePath),
		Changelog: version.Changelog,
		CreatedAt: version.CreatedAt,
	}
}

func toTemplateItem(item model.TaskTemplate) dto.TaskTemplateItem {
	return dto.TaskTemplateItem{
		ID:          item.ID,
		Name:        item.Name,
		Summary:     item.Summary,
		Description: item.Description,
		Category:    item.Category,
		Scene:       item.Scene,
		Guide:       item.Guide,
		ResourceRef: support.SplitCSV(item.ResourceRef),
		UsageCount:  item.UsageCount,
		Status:      item.Status,
		UpdatedAt:   item.UpdatedAt,
	}
}

func toApplicationCase(item model.ApplicationCase) dto.ApplicationCaseItem {
	return dto.ApplicationCaseItem{
		ID:         item.ID,
		Title:      item.Title,
		Summary:    item.Summary,
		Category:   item.Category,
		Guide:      item.Guide,
		CoverImage: item.CoverImage,
		Status:     item.Status,
		UpdatedAt:  item.UpdatedAt,
	}
}

func toDocumentItem(item model.Document) dto.DocumentItem {
	return dto.DocumentItem{
		ID:         item.ID,
		CategoryID: item.CategoryID,
		Category:   item.Category.Name,
		Title:      item.Title,
		Summary:    item.Summary,
		Content:    item.Content,
		DocType:    item.DocType,
		Status:     item.Status,
		UpdatedAt:  item.UpdatedAt,
	}
}

func toVideoTutorialItem(item model.VideoTutorial) dto.VideoTutorialItem {
	return dto.VideoTutorialItem{
		ID:        item.ID,
		Title:     item.Title,
		Summary:   item.Summary,
		Link:      item.Link,
		Category:  item.Category,
		SortOrder: item.SortOrder,
		Active:    item.Active,
		UpdatedAt: item.UpdatedAt,
	}
}

func toAgreementTemplateItem(item model.AgreementTemplate) dto.AgreementTemplateItem {
	return dto.AgreementTemplateItem{
		ID:        item.ID,
		Name:      item.Name,
		Content:   item.Content,
		SortOrder: item.SortOrder,
		Active:    item.Active,
		UpdatedAt: item.UpdatedAt,
	}
}

func toDatasetPrivacyOptionItem(item model.DatasetPrivacyOption) dto.DatasetPrivacyOptionItem {
	return dto.DatasetPrivacyOptionItem{
		ID:          item.ID,
		Code:        item.Code,
		Name:        item.Name,
		Description: item.Description,
		SortOrder:   item.SortOrder,
		Active:      item.Active,
		UpdatedAt:   item.UpdatedAt,
	}
}
