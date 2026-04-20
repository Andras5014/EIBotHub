package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/middleware"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

func (h *Handler) listConversations(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.collaboration.Conversations(claims.UserID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) listConversationMessages(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.collaboration.Messages(claims.UserID, parseUintParam(c, "id"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) sendMessage(c *gin.Context) {
	var input dto.MessageSendRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, messageLabels)))
		return
	}
	claims := middleware.MustClaims(c)
	data, err := h.collaboration.SendDirectMessage(claims.UserID, input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) listWorkspaces(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.collaboration.Workspaces(claims.UserID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) workspaceDetail(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.collaboration.WorkspaceDetail(claims.UserID, parseUintParam(c, "id"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) createWorkspace(c *gin.Context) {
	var input dto.WorkspaceCreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, workspaceLabels)))
		return
	}
	claims := middleware.MustClaims(c)
	data, err := h.collaboration.CreateWorkspace(claims.UserID, input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) addWorkspaceMember(c *gin.Context) {
	var input dto.WorkspaceMemberRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", "成员用户不能为空"))
		return
	}
	claims := middleware.MustClaims(c)
	if err := h.collaboration.AddWorkspaceMember(parseUintParam(c, "id"), claims.UserID, input.UserID); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "member added", nil)
}

func (h *Handler) sendWorkspaceMessage(c *gin.Context) {
	var input dto.WorkspaceMessageRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, messageLabels)))
		return
	}
	claims := middleware.MustClaims(c)
	data, err := h.collaboration.SendWorkspaceMessage(parseUintParam(c, "id"), claims.UserID, input.Content)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) adminConversations(c *gin.Context) {
	data, err := h.collaboration.AdminConversations()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminBlockConversation(c *gin.Context) {
	var input dto.AdminConversationModerationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", "封禁原因格式不正确"))
		return
	}
	claims := middleware.MustClaims(c)
	if err := h.collaboration.AdminBlockConversation(parseUintParam(c, "id"), claims.UserID, input.Reason); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "conversation blocked", nil)
}

func (h *Handler) adminUnblockConversation(c *gin.Context) {
	claims := middleware.MustClaims(c)
	if err := h.collaboration.AdminUnblockConversation(parseUintParam(c, "id"), claims.UserID); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "conversation unblocked", nil)
}

func (h *Handler) adminWorkspaces(c *gin.Context) {
	data, err := h.collaboration.AdminWorkspaces()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminBlockWorkspace(c *gin.Context) {
	var input dto.AdminWorkspaceModerationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", "封禁原因格式不正确"))
		return
	}
	claims := middleware.MustClaims(c)
	if err := h.collaboration.AdminBlockWorkspace(parseUintParam(c, "id"), claims.UserID, input.Reason); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "workspace blocked", nil)
}

func (h *Handler) adminUnblockWorkspace(c *gin.Context) {
	claims := middleware.MustClaims(c)
	if err := h.collaboration.AdminUnblockWorkspace(parseUintParam(c, "id"), claims.UserID); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "workspace unblocked", nil)
}

func (h *Handler) adminRemoveWorkspaceMember(c *gin.Context) {
	var input dto.AdminWorkspaceMemberRemovalRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", "成员用户不能为空"))
		return
	}
	claims := middleware.MustClaims(c)
	if err := h.collaboration.AdminRemoveWorkspaceMember(parseUintParam(c, "id"), claims.UserID, input.UserID); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "workspace member removed", nil)
}
