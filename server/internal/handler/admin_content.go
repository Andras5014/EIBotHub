package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/middleware"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

func (h *Handler) listVideos(c *gin.Context) {
	data, err := h.catalog.Videos()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminPortalModules(c *gin.Context) {
	data, err := h.admin.ModuleSettings()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminUpdatePortalModule(c *gin.Context) {
	var input dto.ModuleSettingRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, moduleSettingLabels)))
		return
	}
	if err := h.admin.UpdateModuleSetting(c.Param("key"), input); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "module setting updated", nil)
}

func (h *Handler) adminHomeHeroConfig(c *gin.Context) {
	data, err := h.admin.HomeHeroConfig()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminUpdateHomeHeroConfig(c *gin.Context) {
	var input dto.HomeHeroConfigRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, homeHeroLabels)))
		return
	}
	data, err := h.admin.UpdateHomeHeroConfig(input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminHomeHighlights(c *gin.Context) {
	data, err := h.admin.HomeHighlights()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminCreateHomeHighlight(c *gin.Context) {
	var input dto.HomeHighlightRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, homeHighlightLabels)))
		return
	}
	data, err := h.admin.CreateHomeHighlight(input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) adminUpdateHomeHighlight(c *gin.Context) {
	var input dto.HomeHighlightRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, homeHighlightLabels)))
		return
	}
	data, err := h.admin.UpdateHomeHighlight(parseUintParam(c, "id"), input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminDeleteHomeHighlight(c *gin.Context) {
	if err := h.admin.DeleteHomeHighlight(parseUintParam(c, "id")); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "home highlight deleted", nil)
}

func (h *Handler) listScenePages(c *gin.Context) {
	data, err := h.portal.ScenePages()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) scenePageDetail(c *gin.Context) {
	data, err := h.portal.ScenePageDetail(c.Param("slug"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminScenePages(c *gin.Context) {
	data, err := h.admin.ScenePages()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminCreateScenePage(c *gin.Context) {
	var input dto.ScenePageRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, scenePageLabels)))
		return
	}
	data, err := h.admin.CreateScenePage(input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) adminUpdateScenePage(c *gin.Context) {
	var input dto.ScenePageRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, scenePageLabels)))
		return
	}
	data, err := h.admin.UpdateScenePage(parseUintParam(c, "id"), input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminDeleteScenePage(c *gin.Context) {
	if err := h.admin.DeleteScenePage(parseUintParam(c, "id")); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "scene page deleted", nil)
}

func (h *Handler) adminRankingConfig(c *gin.Context) {
	data, err := h.admin.RankingConfig()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminUpdateRankingConfig(c *gin.Context) {
	var input dto.RankingConfigRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, rankingConfigLabels)))
		return
	}
	data, err := h.admin.UpdateRankingConfig(input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminFeaturedResources(c *gin.Context) {
	data, err := h.admin.FeaturedResources()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminCreateFeaturedResource(c *gin.Context) {
	var input dto.FeaturedResourceRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, featuredResourceLabels)))
		return
	}
	data, err := h.admin.CreateFeaturedResource(input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) adminUpdateFeaturedResource(c *gin.Context) {
	var input dto.FeaturedResourceRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, featuredResourceLabels)))
		return
	}
	data, err := h.admin.UpdateFeaturedResource(parseUintParam(c, "id"), input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminDeleteFeaturedResource(c *gin.Context) {
	if err := h.admin.DeleteFeaturedResource(parseUintParam(c, "id")); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "featured resource deleted", nil)
}

func (h *Handler) adminTemplates(c *gin.Context) {
	data, err := h.admin.Templates()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminCreateTemplate(c *gin.Context) {
	var input dto.TaskTemplateAdminRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, templateAdminLabels)))
		return
	}
	data, err := h.admin.CreateTemplate(input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) adminUpdateTemplate(c *gin.Context) {
	var input dto.TaskTemplateAdminRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, templateAdminLabels)))
		return
	}
	data, err := h.admin.UpdateTemplate(parseUintParam(c, "id"), input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminDeleteTemplate(c *gin.Context) {
	claims := middleware.MustClaims(c)
	if err := h.admin.DeleteTemplate(parseUintParam(c, "id"), claims.UserID); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "template deleted", nil)
}

func (h *Handler) adminBatchTemplateStatus(c *gin.Context) {
	var input dto.BatchStatusRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, batchStatusLabels)))
		return
	}
	claims := middleware.MustClaims(c)
	if err := h.admin.BatchUpdateTemplateStatus(claims.UserID, input); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "template statuses updated", nil)
}

func (h *Handler) adminBatchDeleteTemplates(c *gin.Context) {
	var input dto.BatchDeleteRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, batchDeleteLabels)))
		return
	}
	claims := middleware.MustClaims(c)
	if err := h.admin.BatchDeleteTemplates(claims.UserID, input); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "templates deleted", nil)
}

func (h *Handler) adminApplicationCases(c *gin.Context) {
	data, err := h.admin.ApplicationCases()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminCreateApplicationCase(c *gin.Context) {
	var input dto.ApplicationCaseAdminRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, applicationCaseAdminLabels)))
		return
	}
	data, err := h.admin.CreateApplicationCase(input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) adminUpdateApplicationCase(c *gin.Context) {
	var input dto.ApplicationCaseAdminRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, applicationCaseAdminLabels)))
		return
	}
	data, err := h.admin.UpdateApplicationCase(parseUintParam(c, "id"), input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminDeleteApplicationCase(c *gin.Context) {
	claims := middleware.MustClaims(c)
	if err := h.admin.DeleteApplicationCase(parseUintParam(c, "id"), claims.UserID); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "application case deleted", nil)
}

func (h *Handler) adminBatchApplicationCaseStatus(c *gin.Context) {
	var input dto.BatchStatusRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, batchStatusLabels)))
		return
	}
	claims := middleware.MustClaims(c)
	if err := h.admin.BatchUpdateApplicationCaseStatus(claims.UserID, input); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "application case statuses updated", nil)
}

func (h *Handler) adminBatchDeleteApplicationCases(c *gin.Context) {
	var input dto.BatchDeleteRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, batchDeleteLabels)))
		return
	}
	claims := middleware.MustClaims(c)
	if err := h.admin.BatchDeleteApplicationCases(claims.UserID, input); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "application cases deleted", nil)
}

func (h *Handler) adminDocCategories(c *gin.Context) {
	data, err := h.admin.DocCategories()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminCreateDocCategory(c *gin.Context) {
	var input dto.DocumentCategoryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, documentCategoryLabels)))
		return
	}
	data, err := h.admin.CreateDocCategory(input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) adminUpdateDocCategory(c *gin.Context) {
	var input dto.DocumentCategoryRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, documentCategoryLabels)))
		return
	}
	data, err := h.admin.UpdateDocCategory(parseUintParam(c, "id"), input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminDeleteDocCategory(c *gin.Context) {
	claims := middleware.MustClaims(c)
	if err := h.admin.DeleteDocCategory(parseUintParam(c, "id"), claims.UserID); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "doc category deleted", nil)
}

func (h *Handler) adminDocuments(c *gin.Context) {
	data, err := h.admin.Documents()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminCreateDocument(c *gin.Context) {
	var input dto.DocumentAdminRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, documentAdminLabels)))
		return
	}
	data, err := h.admin.CreateDocument(input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) adminUpdateDocument(c *gin.Context) {
	var input dto.DocumentAdminRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, documentAdminLabels)))
		return
	}
	data, err := h.admin.UpdateDocument(parseUintParam(c, "id"), input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminDeleteDocument(c *gin.Context) {
	claims := middleware.MustClaims(c)
	if err := h.admin.DeleteDocument(parseUintParam(c, "id"), claims.UserID); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "document deleted", nil)
}

func (h *Handler) adminBatchDocumentStatus(c *gin.Context) {
	var input dto.BatchStatusRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, batchStatusLabels)))
		return
	}
	claims := middleware.MustClaims(c)
	if err := h.admin.BatchUpdateDocumentStatus(claims.UserID, input); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "document statuses updated", nil)
}

func (h *Handler) adminBatchDeleteDocuments(c *gin.Context) {
	var input dto.BatchDeleteRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, batchDeleteLabels)))
		return
	}
	claims := middleware.MustClaims(c)
	if err := h.admin.BatchDeleteDocuments(claims.UserID, input); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "documents deleted", nil)
}

func (h *Handler) adminFAQs(c *gin.Context) {
	data, err := h.admin.FAQs()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminCreateFAQ(c *gin.Context) {
	var input dto.FAQRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, faqAdminLabels)))
		return
	}
	data, err := h.admin.CreateFAQ(input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) adminUpdateFAQ(c *gin.Context) {
	var input dto.FAQRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, faqAdminLabels)))
		return
	}
	data, err := h.admin.UpdateFAQ(parseUintParam(c, "id"), input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminDeleteFAQ(c *gin.Context) {
	claims := middleware.MustClaims(c)
	if err := h.admin.DeleteFAQ(parseUintParam(c, "id"), claims.UserID); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "faq deleted", nil)
}

func (h *Handler) adminVideos(c *gin.Context) {
	data, err := h.admin.VideoTutorials()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminCreateVideo(c *gin.Context) {
	var input dto.VideoTutorialRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, videoTutorialLabels)))
		return
	}
	data, err := h.admin.CreateVideoTutorial(input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) adminUpdateVideo(c *gin.Context) {
	var input dto.VideoTutorialRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, videoTutorialLabels)))
		return
	}
	data, err := h.admin.UpdateVideoTutorial(parseUintParam(c, "id"), input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminDeleteVideo(c *gin.Context) {
	claims := middleware.MustClaims(c)
	if err := h.admin.DeleteVideoTutorial(parseUintParam(c, "id"), claims.UserID); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "video deleted", nil)
}

func (h *Handler) adminBatchVideoStatus(c *gin.Context) {
	var input struct {
		IDs    []uint `json:"ids" binding:"required,min=1,max=100"`
		Active bool   `json:"active"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", "请提供视频列表"))
		return
	}
	claims := middleware.MustClaims(c)
	if err := h.admin.BatchUpdateVideoStatus(claims.UserID, input.IDs, input.Active); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "video statuses updated", nil)
}

func (h *Handler) adminBatchDeleteVideos(c *gin.Context) {
	var input dto.BatchDeleteRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, batchDeleteLabels)))
		return
	}
	claims := middleware.MustClaims(c)
	if err := h.admin.BatchDeleteVideos(claims.UserID, input); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "videos deleted", nil)
}

func (h *Handler) adminOperationLogs(c *gin.Context) {
	data, err := h.admin.OperationLogs()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminAgreementTemplates(c *gin.Context) {
	data, err := h.admin.AgreementTemplates()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminCreateAgreementTemplate(c *gin.Context) {
	var input dto.AgreementTemplateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, agreementTemplateLabels)))
		return
	}
	data, err := h.admin.CreateAgreementTemplate(input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) adminUpdateAgreementTemplate(c *gin.Context) {
	var input dto.AgreementTemplateRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, agreementTemplateLabels)))
		return
	}
	data, err := h.admin.UpdateAgreementTemplate(parseUintParam(c, "id"), input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminDeleteAgreementTemplate(c *gin.Context) {
	if err := h.admin.DeleteAgreementTemplate(parseUintParam(c, "id")); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "agreement template deleted", nil)
}

func (h *Handler) adminPrivacyOptions(c *gin.Context) {
	data, err := h.admin.PrivacyOptions()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminCreatePrivacyOption(c *gin.Context) {
	var input dto.DatasetPrivacyOptionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, privacyOptionLabels)))
		return
	}
	data, err := h.admin.CreatePrivacyOption(input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) adminUpdatePrivacyOption(c *gin.Context) {
	var input dto.DatasetPrivacyOptionRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, privacyOptionLabels)))
		return
	}
	data, err := h.admin.UpdatePrivacyOption(parseUintParam(c, "id"), input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminDeletePrivacyOption(c *gin.Context) {
	if err := h.admin.DeletePrivacyOption(parseUintParam(c, "id")); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "privacy option deleted", nil)
}

func (h *Handler) adminFilterOptions(c *gin.Context) {
	data, err := h.admin.FilterOptionConfigs(c.Query("kind"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminCreateFilterOption(c *gin.Context) {
	var input dto.FilterOptionConfigRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, filterOptionLabels)))
		return
	}
	claims := middleware.MustClaims(c)
	data, err := h.admin.CreateFilterOptionConfig(claims.UserID, input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) adminUpdateFilterOption(c *gin.Context) {
	var input dto.FilterOptionConfigRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, filterOptionLabels)))
		return
	}
	claims := middleware.MustClaims(c)
	data, err := h.admin.UpdateFilterOptionConfig(parseUintParam(c, "id"), claims.UserID, input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminDeleteFilterOption(c *gin.Context) {
	claims := middleware.MustClaims(c)
	if err := h.admin.DeleteFilterOptionConfig(parseUintParam(c, "id"), claims.UserID); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "filter option deleted", nil)
}

func (h *Handler) adminModelRecommendTags(c *gin.Context) {
	data, err := h.models.AdminRecommendTags()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminUpdateModelRecommendTag(c *gin.Context) {
	var input dto.ModelRecommendTagRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, modelRecommendTagLabels)))
		return
	}
	if err := h.models.UpdateRecommendTag(parseUintParam(c, "id"), input); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "model recommend tag updated", nil)
}
