package service

import (
	"errors"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/Andras5014/EIBotHub/server/internal/model"
	"github.com/Andras5014/EIBotHub/server/internal/repository"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

type FileService struct {
	repo    *repository.FileObjectRepository
	storage support.ObjectStorage
	tokens  *support.FileTokenManager
}

type ResolvedFileDownload struct {
	Path         string
	FileName     string
	MIMEType     string
	Inline       bool
	ContentSize  int64
	StorageKey   string
	StorageScope string
}

func NewFileService(repo *repository.FileObjectRepository, storage support.ObjectStorage, tokens *support.FileTokenManager) *FileService {
	return &FileService{repo: repo, storage: storage, tokens: tokens}
}

func (s *FileService) AuthorizedURL(objectKey, originalName string, uploadedBy uint, inline bool, ttl time.Duration) string {
	if objectKey == "" || s == nil || s.repo == nil || s.storage == nil || s.tokens == nil {
		return buildStorageURL(objectKey)
	}
	item, err := s.ensureObject(objectKey, originalName, uploadedBy)
	if err != nil {
		return buildStorageURL(objectKey)
	}
	return "/api/v1/files/download/" + s.tokens.Issue(item.ID, inline, ttl)
}

func (s *FileService) ResolveDownload(token string) (*ResolvedFileDownload, error) {
	if s == nil || s.repo == nil || s.storage == nil || s.tokens == nil {
		return nil, support.NewError(http.StatusNotFound, "file_service_unavailable", "file service unavailable")
	}
	claims, err := s.tokens.Parse(token)
	if err != nil {
		return nil, err
	}
	item, err := s.repo.GetByID(claims.FileID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, support.NewError(http.StatusNotFound, "file_not_found", "file not found")
		}
		return nil, err
	}
	fullPath := s.storage.ResolvePath(item.ObjectKey)
	info, statErr := os.Stat(fullPath)
	if statErr != nil {
		return nil, support.NewError(http.StatusNotFound, "file_not_found", "file not found")
	}
	return &ResolvedFileDownload{
		Path:         fullPath,
		FileName:     item.OriginalName,
		MIMEType:     item.MIMEType,
		Inline:       claims.Inline,
		ContentSize:  info.Size(),
		StorageKey:   item.ObjectKey,
		StorageScope: item.Scope,
	}, nil
}

func (s *FileService) ensureObject(objectKey, originalName string, uploadedBy uint) (*model.FileObject, error) {
	item, err := s.repo.FindByObjectKey(objectKey)
	if err == nil {
		needsUpdate := false
		if item.OriginalName == "" && originalName != "" {
			item.OriginalName = originalName
			needsUpdate = true
		}
		if item.UploadedBy == 0 && uploadedBy > 0 {
			item.UploadedBy = uploadedBy
			needsUpdate = true
		}
		if needsUpdate {
			if updateErr := s.repo.Update(item); updateErr != nil {
				return nil, updateErr
			}
		}
		return item, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	fullPath := s.storage.ResolvePath(objectKey)
	info, statErr := os.Stat(fullPath)
	if statErr != nil {
		return nil, statErr
	}

	if originalName == "" {
		originalName = filepath.Base(objectKey)
	}
	mimeType := mime.TypeByExtension(strings.ToLower(filepath.Ext(originalName)))
	item = &model.FileObject{
		ObjectKey:     objectKey,
		OriginalName:  originalName,
		MIMEType:      mimeType,
		SizeBytes:     info.Size(),
		StorageDriver: "local",
		Scope:         strings.TrimSpace(filepath.ToSlash(filepath.Dir(objectKey))),
		UploadedBy:    uploadedBy,
	}
	if err := s.repo.Create(item); err != nil {
		return nil, err
	}
	return item, nil
}
