package repository

import (
	"fmt"
	"sort"

	"gorm.io/gorm"

	"github.com/Andras5014/EIBotHub/server/internal/model"
)

type CollaborationRepository struct {
	db *gorm.DB
}

func NewCollaborationRepository(db *gorm.DB) *CollaborationRepository {
	return &CollaborationRepository{db: db}
}

func DirectKey(userA, userB uint) string {
	ids := []uint{userA, userB}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return fmt.Sprintf("direct:%d:%d", ids[0], ids[1])
}

func WorkspaceKey(workspaceID uint) string {
	return fmt.Sprintf("workspace:%d", workspaceID)
}

func (r *CollaborationRepository) FindConversationByKey(key string) (*model.Conversation, error) {
	var item model.Conversation
	err := r.db.Where("key = ?", key).First(&item).Error
	return &item, err
}

func (r *CollaborationRepository) CreateConversation(item *model.Conversation) error {
	return r.db.Create(item).Error
}

func (r *CollaborationRepository) UpdateConversation(item *model.Conversation) error {
	return r.db.Save(item).Error
}

func (r *CollaborationRepository) AddParticipants(items []model.ConversationParticipant) error {
	return r.db.Create(&items).Error
}

func (r *CollaborationRepository) ListParticipants(conversationID uint) ([]model.ConversationParticipant, error) {
	var items []model.ConversationParticipant
	err := r.db.Where("conversation_id = ?", conversationID).Order("created_at asc").Find(&items).Error
	return items, err
}

func (r *CollaborationRepository) HasParticipant(conversationID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.ConversationParticipant{}).Where("conversation_id = ? AND user_id = ?", conversationID, userID).Count(&count).Error
	return count > 0, err
}

func (r *CollaborationRepository) ListConversationsByUser(userID uint) ([]model.Conversation, error) {
	var items []model.Conversation
	err := r.db.Model(&model.Conversation{}).
		Joins("join conversation_participants on conversation_participants.conversation_id = conversations.id").
		Where("conversation_participants.user_id = ?", userID).
		Order("conversations.updated_at desc").
		Find(&items).Error
	return items, err
}

func (r *CollaborationRepository) AddMessage(item *model.Message) error {
	if err := r.db.Create(item).Error; err != nil {
		return err
	}
	return r.db.Model(&model.Conversation{}).Where("id = ?", item.ConversationID).Update("updated_at", item.CreatedAt).Error
}

func (r *CollaborationRepository) ListMessages(conversationID uint) ([]model.Message, error) {
	var items []model.Message
	err := r.db.Where("conversation_id = ?", conversationID).Order("created_at asc").Find(&items).Error
	return items, err
}

func (r *CollaborationRepository) LatestMessage(conversationID uint) (*model.Message, error) {
	var item model.Message
	err := r.db.Where("conversation_id = ?", conversationID).Order("created_at desc").First(&item).Error
	return &item, err
}

func (r *CollaborationRepository) CreateWorkspace(item *model.Workspace) error {
	return r.db.Create(item).Error
}

func (r *CollaborationRepository) UpdateWorkspace(item *model.Workspace) error {
	return r.db.Save(item).Error
}

func (r *CollaborationRepository) ListWorkspacesByUser(userID uint) ([]model.Workspace, error) {
	var items []model.Workspace
	err := r.db.Model(&model.Workspace{}).
		Joins("join workspace_members on workspace_members.workspace_id = workspaces.id").
		Where("workspace_members.user_id = ?", userID).
		Order("workspaces.updated_at desc").
		Find(&items).Error
	return items, err
}

func (r *CollaborationRepository) GetWorkspace(id uint) (*model.Workspace, error) {
	var item model.Workspace
	err := r.db.First(&item, id).Error
	return &item, err
}

func (r *CollaborationRepository) AddWorkspaceMember(item *model.WorkspaceMember) error {
	return r.db.Where("workspace_id = ? and user_id = ?", item.WorkspaceID, item.UserID).FirstOrCreate(item).Error
}

func (r *CollaborationRepository) ListWorkspaceMembers(workspaceID uint) ([]model.WorkspaceMember, error) {
	var items []model.WorkspaceMember
	err := r.db.Where("workspace_id = ?", workspaceID).Order("created_at asc").Find(&items).Error
	return items, err
}

func (r *CollaborationRepository) HasWorkspaceMember(workspaceID, userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.WorkspaceMember{}).Where("workspace_id = ? AND user_id = ?", workspaceID, userID).Count(&count).Error
	return count > 0, err
}

func (r *CollaborationRepository) CountWorkspaceMembers(workspaceID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.WorkspaceMember{}).Where("workspace_id = ?", workspaceID).Count(&count).Error
	return count, err
}
