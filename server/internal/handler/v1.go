package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/middleware"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

func (h *Handler) hotQueries(c *gin.Context) {
	data, err := h.community.HotQueries()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) listModelEvaluations(c *gin.Context) {
	data, err := h.community.ListModelEvaluations(parseUintParam(c, "id"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) createModelEvaluation(c *gin.Context) {
	var input dto.ModelEvaluationRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, modelEvaluationLabels)))
		return
	}
	claims := middleware.MustClaims(c)
	data, err := h.community.CreateModelEvaluation(claims.UserID, parseUintParam(c, "id"), input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) getModelRatings(c *gin.Context) {
	h.getRatings(c, "model")
}

func (h *Handler) rateModel(c *gin.Context) {
	h.rateResource(c, "model")
}

func (h *Handler) getModelComments(c *gin.Context) {
	h.getComments(c, "model")
}

func (h *Handler) commentModel(c *gin.Context) {
	h.commentResource(c, "model")
}

func (h *Handler) datasetStats(c *gin.Context) {
	data, err := h.community.DatasetStats(parseUintParam(c, "id"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) getDatasetRatings(c *gin.Context) {
	h.getRatings(c, "dataset")
}

func (h *Handler) rateDataset(c *gin.Context) {
	h.rateResource(c, "dataset")
}

func (h *Handler) getDatasetComments(c *gin.Context) {
	h.getComments(c, "dataset")
}

func (h *Handler) commentDataset(c *gin.Context) {
	h.commentResource(c, "dataset")
}

func (h *Handler) getTemplateRatings(c *gin.Context) {
	h.getRatings(c, "task-template")
}

func (h *Handler) rateTemplate(c *gin.Context) {
	h.rateResource(c, "task-template")
}

func (h *Handler) getTemplateComments(c *gin.Context) {
	h.getComments(c, "task-template")
}

func (h *Handler) commentTemplate(c *gin.Context) {
	h.commentResource(c, "task-template")
}

func (h *Handler) listSkills(c *gin.Context) {
	data, err := h.community.ListSkills()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) skillDetail(c *gin.Context) {
	data, err := h.community.GetSkill(parseUintParam(c, "id"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) createSkill(c *gin.Context) {
	var input dto.SkillCreateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, skillLabels)))
		return
	}
	claims := middleware.MustClaims(c)
	data, err := h.community.CreateSkill(claims.UserID, input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) forkSkill(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.community.ForkSkill(parseUintParam(c, "id"), claims.UserID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) getSkillRatings(c *gin.Context) {
	h.getRatings(c, "skill")
}

func (h *Handler) rateSkill(c *gin.Context) {
	h.rateResource(c, "skill")
}

func (h *Handler) getSkillComments(c *gin.Context) {
	h.getComments(c, "skill")
}

func (h *Handler) commentSkill(c *gin.Context) {
	h.commentResource(c, "skill")
}

func (h *Handler) listDiscussions(c *gin.Context) {
	data, err := h.community.ListDiscussions()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) discussionDetail(c *gin.Context) {
	data, err := h.community.GetDiscussion(parseUintParam(c, "id"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) createDiscussion(c *gin.Context) {
	var input dto.DiscussionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, discussionLabels)))
		return
	}
	claims := middleware.MustClaims(c)
	data, err := h.community.CreateDiscussion(claims.UserID, input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) getDiscussionComments(c *gin.Context) {
	h.getComments(c, "discussion")
}

func (h *Handler) commentDiscussion(c *gin.Context) {
	h.commentResource(c, "discussion")
}

func (h *Handler) toggleFollow(c *gin.Context) {
	claims := middleware.MustClaims(c)
	if err := h.community.ToggleFollow(claims.UserID, parseUintParam(c, "id")); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "follow state updated", nil)
}

func (h *Handler) myFollows(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.community.MyFollows(claims.UserID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) myFollowers(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.community.MyFollowers(claims.UserID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) myFollowStats(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.community.FollowStats(claims.UserID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) myContributions(c *gin.Context) {
	claims := middleware.MustClaims(c)
	data, err := h.community.UserContributions(claims.UserID)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) publicUserContributions(c *gin.Context) {
	data, err := h.community.UserContributions(parseUintParam(c, "id"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) publicUserFollowStats(c *gin.Context) {
	data, err := h.community.FollowStats(parseUintParam(c, "id"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminCommunityOverview(c *gin.Context) {
	data, err := h.community.AdminOverview()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminCommunitySkills(c *gin.Context) {
	data, err := h.community.AdminSkills()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminHideSkill(c *gin.Context) {
	if err := h.community.AdminHideSkill(parseUintParam(c, "id")); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "skill hidden", nil)
}

func (h *Handler) adminCommunityDiscussions(c *gin.Context) {
	data, err := h.community.AdminDiscussions()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminDeleteDiscussion(c *gin.Context) {
	if err := h.community.AdminDeleteDiscussion(parseUintParam(c, "id")); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "discussion removed", nil)
}

func (h *Handler) adminCommunityComments(c *gin.Context) {
	data, err := h.community.AdminComments()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminDeleteComment(c *gin.Context) {
	if err := h.community.AdminDeleteComment(parseUintParam(c, "id")); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "comment removed", nil)
}

func (h *Handler) getRatings(c *gin.Context, resourceType string) {
	data, err := h.community.RatingSummary(resourceType, parseUintParam(c, "id"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) rateResource(c *gin.Context, resourceType string) {
	var input dto.RatingRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, ratingLabels)))
		return
	}
	claims := middleware.MustClaims(c)
	data, err := h.community.UpsertRating(resourceType, parseUintParam(c, "id"), claims.UserID, input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) getComments(c *gin.Context, resourceType string) {
	data, err := h.community.ListComments(resourceType, parseUintParam(c, "id"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) commentResource(c *gin.Context, resourceType string) {
	var input dto.CommentRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, commentLabels)))
		return
	}
	claims := middleware.MustClaims(c)
	data, err := h.community.AddComment(resourceType, parseUintParam(c, "id"), claims.UserID, input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}
