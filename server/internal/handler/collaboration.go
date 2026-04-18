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
