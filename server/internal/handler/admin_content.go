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
