package service

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/model"
	"github.com/Andras5014/EIBotHub/server/internal/repository"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

type CollaborationService struct {
	repo     *repository.CollaborationRepository
	users    *repository.UserRepository
	activity *repository.UserActivityRepository
	ops      *repository.OperationLogRepository
}

func NewCollaborationService(repo *repository.CollaborationRepository, users *repository.UserRepository, activity *repository.UserActivityRepository, ops *repository.OperationLogRepository) *CollaborationService {
	return &CollaborationService{repo: repo, users: users, activity: activity, ops: ops}
}

func (s *CollaborationService) Conversations(userID uint) ([]dto.ConversationItem, error) {
	items, err := s.repo.ListConversationsByUser(userID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.ConversationItem, 0, len(items))
	for _, item := range items {
		participants, _ := s.repo.ListParticipants(item.ID)
		names := make([]string, 0, len(participants))
		for _, participant := range participants {
			user, _ := s.users.FindByID(participant.UserID)
			names = append(names, user.Username)
		}
		latestMessage, _ := s.repo.LatestMessage(item.ID)
		result = append(result, dto.ConversationItem{
			ID:               item.ID,
			Kind:             item.Kind,
			Title:            item.Title,
			WorkspaceID:      item.WorkspaceID,
			ParticipantNames: names,
			LatestMessage:    messageContent(latestMessage),
			Active:           item.Active,
			BlockedReason:    item.BlockedReason,
			UpdatedAt:        item.UpdatedAt,
		})
	}
	return result, nil
}

func (s *CollaborationService) Messages(userID, conversationID uint) ([]dto.MessageItem, error) {
	conversation, err := s.repo.GetConversation(conversationID)
	if err != nil {
		return nil, support.NewError(http.StatusNotFound, "conversation_not_found", "会话不存在")
	}
	if !conversation.Active {
		return nil, support.NewError(http.StatusForbidden, "conversation_blocked", "当前会话已被封禁")
	}
	allowed, err := s.repo.HasParticipant(conversationID, userID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, support.NewError(http.StatusForbidden, "forbidden", "无权访问该会话")
	}
	items, err := s.repo.ListMessages(conversationID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.MessageItem, 0, len(items))
	for _, item := range items {
		user, _ := s.users.FindByID(item.SenderID)
		result = append(result, dto.MessageItem{
			ID:             item.ID,
			ConversationID: item.ConversationID,
			SenderID:       item.SenderID,
			SenderName:     user.Username,
			Content:        item.Content,
			CreatedAt:      item.CreatedAt,
		})
	}
	return result, nil
}

func (s *CollaborationService) SendDirectMessage(senderID uint, input dto.MessageSendRequest) (*dto.MessageItem, error) {
	if input.ConversationID != 0 {
		conversation, err := s.repo.GetConversation(input.ConversationID)
		if err != nil {
			return nil, support.NewError(http.StatusNotFound, "conversation_not_found", "会话不存在")
		}
		if !conversation.Active {
			return nil, support.NewError(http.StatusForbidden, "conversation_blocked", "当前会话已被封禁")
		}
		allowed, err := s.repo.HasParticipant(input.ConversationID, senderID)
		if err != nil {
			return nil, err
		}
		if !allowed {
			return nil, support.NewError(http.StatusForbidden, "forbidden", "无权访问该会话")
		}
		participants, err := s.repo.ListParticipants(input.ConversationID)
		if err != nil {
			return nil, err
		}
		var notifyUserID uint
		for _, participant := range participants {
			if participant.UserID != senderID {
				notifyUserID = participant.UserID
				break
			}
		}
		return s.addMessage(input.ConversationID, senderID, input.Content, notifyUserID)
	}

	if input.RecipientUserID == 0 {
		return nil, support.NewError(http.StatusBadRequest, "recipient_required", "接收用户不能为空")
	}
	if _, err := s.users.FindByID(input.RecipientUserID); err != nil {
		return nil, support.NewError(http.StatusNotFound, "user_not_found", "接收用户不存在")
	}
	key := repository.DirectKey(senderID, input.RecipientUserID)
	conversation, err := s.repo.FindConversationByKey(key)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		title := "私信会话"
		conversation = &model.Conversation{
			Kind:  "direct",
			Key:   key,
			Title: title,
		}
		if err := s.repo.CreateConversation(conversation); err != nil {
			return nil, err
		}
		if err := s.repo.AddParticipants([]model.ConversationParticipant{
			{ConversationID: conversation.ID, UserID: senderID},
			{ConversationID: conversation.ID, UserID: input.RecipientUserID},
		}); err != nil {
			return nil, err
		}
	} else if !conversation.Active {
		return nil, support.NewError(http.StatusForbidden, "conversation_blocked", "当前会话已被封禁")
	}

	return s.addMessage(conversation.ID, senderID, input.Content, input.RecipientUserID)
}

func (s *CollaborationService) Workspaces(userID uint) ([]dto.WorkspaceItem, error) {
	items, err := s.repo.ListWorkspacesByUser(userID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.WorkspaceItem, 0, len(items))
	for _, item := range items {
		count, _ := s.repo.CountWorkspaceMembers(item.ID)
		result = append(result, dto.WorkspaceItem{
			ID:             item.ID,
			Name:           item.Name,
			Summary:        item.Summary,
			OwnerID:        item.OwnerID,
			ConversationID: item.ConversationID,
			MemberCount:    count,
			Active:         item.Active,
			BlockedReason:  item.BlockedReason,
			UpdatedAt:      item.UpdatedAt,
		})
	}
	return result, nil
}

func (s *CollaborationService) WorkspaceDetail(userID, id uint) (*dto.WorkspaceDetail, error) {
	item, err := s.repo.GetWorkspace(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, support.NewError(http.StatusNotFound, "workspace_not_found", "协作空间不存在")
		}
		return nil, err
	}
	if !item.Active {
		return nil, support.NewError(http.StatusForbidden, "workspace_blocked", "当前协作空间已被封禁")
	}
	allowed, err := s.repo.HasWorkspaceMember(id, userID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, support.NewError(http.StatusForbidden, "forbidden", "无权访问该协作空间")
	}
	memberRows, err := s.repo.ListWorkspaceMembers(id)
	if err != nil {
		return nil, err
	}
	members := make([]dto.FollowItem, 0, len(memberRows))
	for _, member := range memberRows {
		user, _ := s.users.FindByID(member.UserID)
		members = append(members, dto.FollowItem{
			ID:        member.ID,
			UserID:    user.ID,
			UserName:  user.Username,
			Bio:       user.Bio,
			CreatedAt: member.CreatedAt,
		})
	}
	messages, err := s.Messages(userID, item.ConversationID)
	if err != nil {
		return nil, err
	}
	count, _ := s.repo.CountWorkspaceMembers(id)
	return &dto.WorkspaceDetail{
		WorkspaceItem: dto.WorkspaceItem{
			ID:             item.ID,
			Name:           item.Name,
			Summary:        item.Summary,
			OwnerID:        item.OwnerID,
			ConversationID: item.ConversationID,
			MemberCount:    count,
			Active:         item.Active,
			BlockedReason:  item.BlockedReason,
			UpdatedAt:      item.UpdatedAt,
		},
		Messages: messages,
		Members:  members,
	}, nil
}

func (s *CollaborationService) CreateWorkspace(ownerID uint, input dto.WorkspaceCreateRequest) (*dto.WorkspaceDetail, error) {
	conversation := &model.Conversation{
		Kind:  "workspace",
		Key:   "workspace:pending:" + time.Now().Format("20060102150405.000000"),
		Title: input.Name,
	}
	if err := s.repo.CreateConversation(conversation); err != nil {
		return nil, err
	}
	workspace := &model.Workspace{
		Name:           input.Name,
		Summary:        input.Summary,
		OwnerID:        ownerID,
		ConversationID: conversation.ID,
	}
	if err := s.repo.CreateWorkspace(workspace); err != nil {
		return nil, err
	}
	conversation.Key = repository.WorkspaceKey(workspace.ID)
	conversation.WorkspaceID = &workspace.ID
	conversation.Title = workspace.Name
	if err := s.repo.UpdateConversation(conversation); err != nil {
		return nil, err
	}
	if err := s.repo.AddWorkspaceMember(&model.WorkspaceMember{
		WorkspaceID: workspace.ID,
		UserID:      ownerID,
		Role:        "owner",
	}); err != nil {
		return nil, err
	}
	if err := s.repo.AddParticipants([]model.ConversationParticipant{
		{ConversationID: conversation.ID, UserID: ownerID},
	}); err != nil {
		return nil, err
	}
	return s.WorkspaceDetail(ownerID, workspace.ID)
}

func (s *CollaborationService) AddWorkspaceMember(workspaceID uint, actorUserID uint, userID uint) error {
	item, err := s.repo.GetWorkspace(workspaceID)
	if err != nil {
		return support.NewError(http.StatusNotFound, "workspace_not_found", "协作空间不存在")
	}
	if !item.Active {
		return support.NewError(http.StatusForbidden, "workspace_blocked", "当前协作空间已被封禁")
	}
	if item.OwnerID != actorUserID {
		return support.NewError(http.StatusForbidden, "forbidden", "只有空间所有者可以添加成员")
	}
	if _, err := s.users.FindByID(userID); err != nil {
		return support.NewError(http.StatusNotFound, "user_not_found", "成员不存在")
	}
	if err := s.repo.AddWorkspaceMember(&model.WorkspaceMember{
		WorkspaceID: workspaceID,
		UserID:      userID,
		Role:        "member",
	}); err != nil {
		return err
	}
	return s.repo.AddParticipants([]model.ConversationParticipant{{ConversationID: item.ConversationID, UserID: userID}})
}

func (s *CollaborationService) SendWorkspaceMessage(workspaceID, senderID uint, content string) (*dto.MessageItem, error) {
	item, err := s.repo.GetWorkspace(workspaceID)
	if err != nil {
		return nil, support.NewError(http.StatusNotFound, "workspace_not_found", "协作空间不存在")
	}
	if !item.Active {
		return nil, support.NewError(http.StatusForbidden, "workspace_blocked", "当前协作空间已被封禁")
	}
	allowed, err := s.repo.HasWorkspaceMember(workspaceID, senderID)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, support.NewError(http.StatusForbidden, "forbidden", "无权访问该协作空间")
	}
	members, _ := s.repo.ListWorkspaceMembers(workspaceID)
	var notifyUserID uint
	if len(members) > 0 {
		for _, member := range members {
			if member.UserID != senderID {
				notifyUserID = member.UserID
				break
			}
		}
	}
	return s.addMessage(item.ConversationID, senderID, content, notifyUserID)
}

func (s *CollaborationService) addMessage(conversationID, senderID uint, content string, notifyUserID uint) (*dto.MessageItem, error) {
	item := &model.Message{
		ConversationID: conversationID,
		SenderID:       senderID,
		Content:        content,
		CreatedAt:      time.Now(),
	}
	if err := s.repo.AddMessage(item); err != nil {
		return nil, err
	}
	sender, _ := s.users.FindByID(senderID)
	if notifyUserID != 0 && notifyUserID != senderID && s.activity != nil {
		_ = s.activity.AddNotification(&model.Notification{
			UserID:  notifyUserID,
			Type:    "message",
			Title:   "你收到新的消息",
			Content: sender.Username + " 发送了新的消息",
			Link:    "/messages?conversation=" + strconv.FormatUint(uint64(conversationID), 10),
		})
	}
	return &dto.MessageItem{
		ID:             item.ID,
		ConversationID: item.ConversationID,
		SenderID:       item.SenderID,
		SenderName:     sender.Username,
		Content:        item.Content,
		CreatedAt:      item.CreatedAt,
	}, nil
}

func messageContent(item *model.Message) string {
	if item == nil {
		return ""
	}
	return item.Content
}

func (s *CollaborationService) AdminConversations() ([]dto.AdminConversationModerationItem, error) {
	items, err := s.repo.ListAllConversations()
	if err != nil {
		return nil, err
	}
	result := make([]dto.AdminConversationModerationItem, 0, len(items))
	for _, item := range items {
		participants, _ := s.repo.ListParticipants(item.ID)
		names := make([]string, 0, len(participants))
		for _, participant := range participants {
			user, _ := s.users.FindByID(participant.UserID)
			names = append(names, user.Username)
		}
		latestMessage, _ := s.repo.LatestMessage(item.ID)
		result = append(result, dto.AdminConversationModerationItem{
			ID:               item.ID,
			Kind:             item.Kind,
			Title:            item.Title,
			ParticipantNames: names,
			LatestMessage:    messageContent(latestMessage),
			Active:           item.Active,
			BlockedReason:    item.BlockedReason,
			UpdatedAt:        item.UpdatedAt,
		})
	}
	return result, nil
}

func (s *CollaborationService) AdminBlockConversation(id, adminUserID uint, reason string) error {
	item, err := s.repo.GetConversation(id)
	if err != nil {
		return support.NewError(http.StatusNotFound, "conversation_not_found", "会话不存在")
	}
	now := time.Now()
	item.Active = false
	item.BlockedReason = reason
	item.BlockedBy = &adminUserID
	item.BlockedAt = &now
	if err := s.repo.UpdateConversation(item); err != nil {
		return err
	}
	return s.logOperation(adminUserID, "block", "conversation", item.ID, item.Title, reason)
}

func (s *CollaborationService) AdminUnblockConversation(id, adminUserID uint) error {
	item, err := s.repo.GetConversation(id)
	if err != nil {
		return support.NewError(http.StatusNotFound, "conversation_not_found", "会话不存在")
	}
	item.Active = true
	item.BlockedReason = ""
	item.BlockedBy = nil
	item.BlockedAt = nil
	if err := s.repo.UpdateConversation(item); err != nil {
		return err
	}
	return s.logOperation(adminUserID, "unblock", "conversation", item.ID, item.Title, "")
}

func (s *CollaborationService) AdminWorkspaces() ([]dto.AdminWorkspaceModerationItem, error) {
	items, err := s.repo.ListAllWorkspaces()
	if err != nil {
		return nil, err
	}
	result := make([]dto.AdminWorkspaceModerationItem, 0, len(items))
	for _, item := range items {
		memberRows, _ := s.repo.ListWorkspaceMembers(item.ID)
		members := make([]dto.FollowItem, 0, len(memberRows))
		for _, member := range memberRows {
			user, _ := s.users.FindByID(member.UserID)
			members = append(members, dto.FollowItem{
				ID:        member.ID,
				UserID:    user.ID,
				UserName:  user.Username,
				Bio:       user.Bio,
				CreatedAt: member.CreatedAt,
			})
		}
		ownerName := ""
		if owner, ownerErr := s.users.FindByID(item.OwnerID); ownerErr == nil {
			ownerName = owner.Username
		}
		result = append(result, dto.AdminWorkspaceModerationItem{
			ID:            item.ID,
			Name:          item.Name,
			Summary:       item.Summary,
			OwnerID:       item.OwnerID,
			OwnerName:     ownerName,
			MemberCount:   int64(len(members)),
			Members:       members,
			Active:        item.Active,
			BlockedReason: item.BlockedReason,
			UpdatedAt:     item.UpdatedAt,
		})
	}
	return result, nil
}

func (s *CollaborationService) AdminBlockWorkspace(id, adminUserID uint, reason string) error {
	item, err := s.repo.GetWorkspace(id)
	if err != nil {
		return support.NewError(http.StatusNotFound, "workspace_not_found", "协作空间不存在")
	}
	now := time.Now()
	item.Active = false
	item.BlockedReason = reason
	item.BlockedBy = &adminUserID
	item.BlockedAt = &now
	if err := s.repo.UpdateWorkspace(item); err != nil {
		return err
	}
	return s.logOperation(adminUserID, "block", "workspace", item.ID, item.Name, reason)
}

func (s *CollaborationService) AdminUnblockWorkspace(id, adminUserID uint) error {
	item, err := s.repo.GetWorkspace(id)
	if err != nil {
		return support.NewError(http.StatusNotFound, "workspace_not_found", "协作空间不存在")
	}
	item.Active = true
	item.BlockedReason = ""
	item.BlockedBy = nil
	item.BlockedAt = nil
	if err := s.repo.UpdateWorkspace(item); err != nil {
		return err
	}
	return s.logOperation(adminUserID, "unblock", "workspace", item.ID, item.Name, "")
}

func (s *CollaborationService) AdminRemoveWorkspaceMember(workspaceID, adminUserID, userID uint) error {
	item, err := s.repo.GetWorkspace(workspaceID)
	if err != nil {
		return support.NewError(http.StatusNotFound, "workspace_not_found", "协作空间不存在")
	}
	if item.OwnerID == userID {
		return support.NewError(http.StatusBadRequest, "workspace_owner_protected", "不能移除协作空间所有者")
	}
	memberExists, err := s.repo.HasWorkspaceMember(workspaceID, userID)
	if err != nil {
		return err
	}
	if !memberExists {
		return support.NewError(http.StatusNotFound, "workspace_member_not_found", "成员不在该协作空间中")
	}
	if err := s.repo.RemoveWorkspaceMember(workspaceID, userID); err != nil {
		return err
	}
	if err := s.repo.RemoveConversationParticipant(item.ConversationID, userID); err != nil {
		return err
	}
	return s.logOperation(adminUserID, "remove_member", "workspace", item.ID, item.Name, strconv.FormatUint(uint64(userID), 10))
}

func (s *CollaborationService) logOperation(adminUserID uint, action, resourceType string, resourceID uint, summary, detail string) error {
	if s.ops == nil {
		return nil
	}
	return s.ops.Create(&model.AdminOperationLog{
		AdminUserID:  adminUserID,
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		Summary:      summary,
		Detail:       detail,
		CreatedAt:    time.Now(),
	})
}
