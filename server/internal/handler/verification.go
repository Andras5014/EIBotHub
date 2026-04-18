package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/middleware"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

func (h *Handler) applyVerification(c *gin.Context) {
	var input dto.DeveloperVerificationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, verificationLabels)))
		return
	}
	claims := middleware.MustClaims(c)
	data, err := h.verification.Apply(claims.UserID, input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) applyEnterpriseVerification(c *gin.Context) {
	var input dto.EnterpriseVerificationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, verificationLabels)))
		return
	}
	claims := middleware.MustClaims(c)
	data, err := h.verification.ApplyEnterprise(claims.UserID, input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) myVerification(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.verification.MyStatus(claims.UserID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) publicUserVerification(c *gin.Context) {
	data, err := h.verification.PublicStatus(parseUintParam(c, "id"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminVerifications(c *gin.Context) {
	data, err := h.verification.AdminList()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminReviewVerification(c *gin.Context) {
	var input dto.VerificationDecisionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, verificationDecisionLabels)))
		return
	}
	claims := middleware.MustClaims(c)
	if err := h.verification.Review(parseUintParam(c, "id"), claims.UserID, input); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "verification reviewed", nil)
}
