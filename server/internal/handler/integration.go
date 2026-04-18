package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/middleware"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

func (h *Handler) openAPISpec(c *gin.Context) {
	support.RespondOK(c, h.integrations.Spec())
}

func (h *Handler) listWebhooks(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.integrations.ListWebhooks(claims.UserID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) createWebhook(c *gin.Context) {
	claims := middleware.MustClaims(c)
	var input dto.WebhookSubscriptionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, webhookLabels)))
		return
	}
	data, err := h.integrations.CreateWebhook(claims.UserID, input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) listWebhookDeliveries(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.integrations.Deliveries(claims.UserID, parseUintParam(c, "id"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) testWebhook(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.integrations.Test(claims.UserID, parseUintParam(c, "id"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}
