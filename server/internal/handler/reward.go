package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/middleware"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

func (h *Handler) myRewardSummary(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.rewards.Summary(claims.UserID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) myRewardLedger(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.rewards.Ledger(claims.UserID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) contributorRankings(c *gin.Context) {
	data, err := h.rewards.Rankings()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) rewardBenefits(c *gin.Context) {
	data, err := h.rewards.Benefits()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) myRewardRedemptions(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.rewards.Redemptions(claims.UserID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) redeemRewardBenefit(c *gin.Context) {
	claims := middleware.MustClaims(c)
	var input dto.RewardRedeemRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, rewardRedeemLabels)))
		return
	}
	data, err := h.rewards.Redeem(claims.UserID, input.BenefitID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	h.rewards.EmitRedeemed(claims.UserID, *data)
	support.RespondCreated(c, data)
}

func (h *Handler) adminRewardOverview(c *gin.Context) {
	data, err := h.rewards.AdminOverview()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminRewardBenefits(c *gin.Context) {
	data, err := h.rewards.AdminBenefits()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminUpdateRewardBenefit(c *gin.Context) {
	var input dto.AdminRewardBenefitRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, adminRewardBenefitLabels)))
		return
	}
	data, err := h.rewards.UpdateBenefit(parseUintParam(c, "id"), input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminRewardAdjustments(c *gin.Context) {
	data, err := h.rewards.AdminAdjustments()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminAdjustRewardPoints(c *gin.Context) {
	var input dto.AdminRewardAdjustmentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, adminRewardAdjustmentLabels)))
		return
	}
	data, err := h.rewards.AdjustPoints(input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}
