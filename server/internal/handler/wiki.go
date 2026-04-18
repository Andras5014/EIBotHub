package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Andras5014/EIBotHub/server/internal/dto"
	"github.com/Andras5014/EIBotHub/server/internal/middleware"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

func (h *Handler) listWikiPages(c *gin.Context) {
	data, err := h.wiki.ListPages()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) wikiPageDetail(c *gin.Context) {
	data, err := h.wiki.GetPage(parseUintParam(c, "id"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) createWikiPage(c *gin.Context) {
	var input dto.WikiPageRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, map[string]string{
			"Title":   "标题",
			"Summary": "摘要",
			"Content": "内容",
			"Comment": "修订说明",
		})))
		return
	}
	claims := middleware.MustClaims(c)
	data, err := h.wiki.CreatePage(claims.UserID, input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondCreated(c, data)
}

func (h *Handler) updateWikiPage(c *gin.Context) {
	var input dto.WikiPageRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, map[string]string{
			"Title":   "标题",
			"Summary": "摘要",
			"Content": "内容",
			"Comment": "修订说明",
		})))
		return
	}
	claims := middleware.MustClaims(c)
	data, err := h.wiki.UpdatePage(parseUintParam(c, "id"), claims.UserID, input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) listWikiRevisions(c *gin.Context) {
	data, err := h.wiki.Revisions(parseUintParam(c, "id"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminWikiPages(c *gin.Context) {
	data, err := h.wiki.AdminPages()
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}

func (h *Handler) adminLockWikiPage(c *gin.Context) {
	claims := middleware.MustClaims(c)
	if err := h.wiki.SetLocked(parseUintParam(c, "id"), claims.UserID, true); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "wiki locked", nil)
}

func (h *Handler) adminUnlockWikiPage(c *gin.Context) {
	claims := middleware.MustClaims(c)
	if err := h.wiki.SetLocked(parseUintParam(c, "id"), claims.UserID, false); err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondMessage(c, "wiki unlocked", nil)
}

func (h *Handler) adminRollbackWikiPage(c *gin.Context) {
	claims := middleware.MustClaims(c)
	var input dto.AdminWikiRollbackRequest
	if err := c.ShouldBindJSON(&input); err != nil {
		support.RespondError(c, support.NewError(http.StatusBadRequest, "invalid_request", support.ValidationMessage(err, wikiRollbackLabels)))
		return
	}
	data, err := h.wiki.Rollback(parseUintParam(c, "id"), claims.UserID, input)
	if err != nil {
		support.RespondError(c, err)
		return
	}
	support.RespondOK(c, data)
}
