package service

import (
	"errors"
	"net/http"
	"time"

	"gorm.io/gorm"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/model"
	"github.com/Andras5014/EIBotHub/server/internal/repository"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

type WikiService struct {
	repo  *repository.WikiRepository
	users *repository.UserRepository
	hooks *IntegrationService
}

func NewWikiService(repo *repository.WikiRepository, users *repository.UserRepository, hooks *IntegrationService) *WikiService {
	return &WikiService{repo: repo, users: users, hooks: hooks}
}

func (s *WikiService) ListPages() ([]dto.WikiPageItem, error) {
	items, err := s.repo.ListPages()
	if err != nil {
		return nil, err
	}
	result := make([]dto.WikiPageItem, 0, len(items))
	for _, item := range items {
		pageItem, err := s.toWikiPageItem(item)
		if err != nil {
			return nil, err
		}
		result = append(result, pageItem)
	}
	return result, nil
}

func (s *WikiService) GetPage(id uint) (*dto.WikiPageItem, error) {
	item, err := s.repo.GetPage(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, support.NewError(http.StatusNotFound, "wiki_not_found", "Wiki 词条不存在")
		}
		return nil, err
	}
	pageItem, err := s.toWikiPageItem(*item)
	if err != nil {
		return nil, err
	}
	result := &pageItem
	return result, nil
}

func (s *WikiService) CreatePage(userID uint, input dto.WikiPageRequest) (*dto.WikiPageItem, error) {
	item := &model.WikiPage{
		Title:    input.Title,
		Summary:  input.Summary,
		Content:  input.Content,
		Status:   model.StatusPublished,
		EditorID: userID,
	}
	if err := s.repo.CreatePage(item); err != nil {
		return nil, err
	}
	if err := s.repo.AddRevision(&model.WikiRevision{
		PageID:   item.ID,
		EditorID: userID,
		Title:    input.Title,
		Summary:  input.Summary,
		Content:  input.Content,
		Comment:  input.Comment,
	}); err != nil {
		return nil, err
	}
	result, err := s.GetPage(item.ID)
	if err != nil {
		return nil, err
	}
	s.emitPageEvent(userID, result, "created")
	return result, nil
}

func (s *WikiService) UpdatePage(id, userID uint, input dto.WikiPageRequest) (*dto.WikiPageItem, error) {
	item, err := s.repo.GetPage(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, support.NewError(http.StatusNotFound, "wiki_not_found", "Wiki 词条不存在")
		}
		return nil, err
	}
	if err := s.ensureEditable(*item, userID); err != nil {
		return nil, err
	}
	item.Title = input.Title
	item.Summary = input.Summary
	item.Content = input.Content
	item.EditorID = userID
	if err := s.repo.UpdatePage(item); err != nil {
		return nil, err
	}
	if err := s.repo.AddRevision(&model.WikiRevision{
		PageID:   item.ID,
		EditorID: userID,
		Title:    input.Title,
		Summary:  input.Summary,
		Content:  input.Content,
		Comment:  input.Comment,
	}); err != nil {
		return nil, err
	}
	result, err := s.GetPage(item.ID)
	if err != nil {
		return nil, err
	}
	s.emitPageEvent(userID, result, "updated")
	return result, nil
}

func (s *WikiService) Revisions(pageID uint) ([]dto.WikiRevisionItem, error) {
	items, err := s.repo.ListRevisions(pageID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.WikiRevisionItem, 0, len(items))
	for _, item := range items {
		user, _ := s.users.FindByID(item.EditorID)
		result = append(result, dto.WikiRevisionItem{
			ID:         item.ID,
			PageID:     item.PageID,
			EditorID:   item.EditorID,
			EditorName: user.Username,
			Title:      item.Title,
			Summary:    item.Summary,
			Content:    item.Content,
			Comment:    item.Comment,
			CreatedAt:  item.CreatedAt,
		})
	}
	return result, nil
}

func (s *WikiService) AdminPages() ([]dto.WikiPageItem, error) {
	items, err := s.repo.ListAllPages()
	if err != nil {
		return nil, err
	}
	result := make([]dto.WikiPageItem, 0, len(items))
	for _, item := range items {
		pageItem, err := s.toWikiPageItem(item)
		if err != nil {
			return nil, err
		}
		result = append(result, pageItem)
	}
	return result, nil
}

func (s *WikiService) SetLocked(pageID, adminUserID uint, locked bool) error {
	item, err := s.repo.GetPage(pageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return support.NewError(http.StatusNotFound, "wiki_not_found", "Wiki 词条不存在")
		}
		return err
	}
	item.Locked = locked
	if locked {
		now := time.Now()
		item.LockedBy = &adminUserID
		item.LockedAt = &now
	} else {
		item.LockedBy = nil
		item.LockedAt = nil
	}
	return s.repo.UpdatePage(item)
}

func (s *WikiService) Rollback(pageID, adminUserID uint, input dto.AdminWikiRollbackRequest) (*dto.WikiPageItem, error) {
	item, err := s.repo.GetPage(pageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, support.NewError(http.StatusNotFound, "wiki_not_found", "Wiki 词条不存在")
		}
		return nil, err
	}
	revision, err := s.repo.GetRevision(pageID, input.RevisionID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, support.NewError(http.StatusNotFound, "wiki_revision_not_found", "Wiki 修订不存在")
		}
		return nil, err
	}
	item.Title = revision.Title
	item.Summary = revision.Summary
	item.Content = revision.Content
	item.EditorID = adminUserID
	if err := s.repo.UpdatePage(item); err != nil {
		return nil, err
	}
	comment := input.Comment
	if comment == "" {
		comment = "管理员回滚到修订版本"
	}
	if err := s.repo.AddRevision(&model.WikiRevision{
		PageID:   item.ID,
		EditorID: adminUserID,
		Title:    item.Title,
		Summary:  item.Summary,
		Content:  item.Content,
		Comment:  comment,
	}); err != nil {
		return nil, err
	}
	result, err := s.GetPage(item.ID)
	if err != nil {
		return nil, err
	}
	s.emitPageEvent(adminUserID, result, "rolled_back")
	return result, nil
}

func (s *WikiService) toWikiPageItem(item model.WikiPage) (dto.WikiPageItem, error) {
	revisionCount, err := s.repo.CountRevisions(item.ID)
	if err != nil {
		return dto.WikiPageItem{}, err
	}
	lockedByName := ""
	if item.LockedBy != nil {
		if user, userErr := s.users.FindByID(*item.LockedBy); userErr == nil {
			lockedByName = user.Username
		}
	}
	return dto.WikiPageItem{
		ID:            item.ID,
		Title:         item.Title,
		Summary:       item.Summary,
		Content:       item.Content,
		Status:        item.Status,
		Locked:        item.Locked,
		LockedBy:      item.LockedBy,
		LockedByName:  lockedByName,
		LockedAt:      item.LockedAt,
		RevisionCount: revisionCount,
		EditorID:      item.EditorID,
		EditorName:    item.Editor.Username,
		UpdatedAt:     item.UpdatedAt,
	}, nil
}

func (s *WikiService) ensureEditable(item model.WikiPage, userID uint) error {
	if !item.Locked {
		return nil
	}
	user, err := s.users.FindByID(userID)
	if err != nil {
		return err
	}
	if model.IsAdminRole(user.Role) {
		return nil
	}
	return support.NewError(http.StatusLocked, "wiki_locked", "该 Wiki 词条已被锁定，仅管理员可编辑")
}

func (s *WikiService) emitPageEvent(userID uint, page *dto.WikiPageItem, action string) {
	if s == nil || s.hooks == nil || page == nil {
		return
	}
	s.hooks.Emit(userID, WebhookEventWikiUpdated, map[string]any{
		"action":      action,
		"page_id":     page.ID,
		"title":       page.Title,
		"editor_id":   page.EditorID,
		"editor_name": page.EditorName,
		"updated_at":  page.UpdatedAt,
	})
}
