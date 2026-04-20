package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Andras5014/EIBotHub/server/internal/support"
)

func (h *Handler) downloadFileByToken(c *gin.Context) {
	if h.files == nil {
		support.RespondError(c, support.NewError(http.StatusNotFound, "file_service_unavailable", "file service unavailable"))
		return
	}
	item, err := h.files.ResolveDownload(c.Param("token"))
	if err != nil {
		support.RespondError(c, err)
		return
	}
	if item.MIMEType != "" {
		c.Header("Content-Type", item.MIMEType)
	}
	if item.Inline {
		c.Header("Content-Disposition", `inline; filename="`+item.FileName+`"`)
		c.File(item.Path)
		return
	}
	c.FileAttachment(item.Path, item.FileName)
}
