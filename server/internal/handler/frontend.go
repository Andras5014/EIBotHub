package handler

import (
	"io/fs"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Andras5014/EIBotHub/server/internal/frontend"
	"github.com/Andras5014/EIBotHub/server/internal/support"
)

func (h *Handler) registerFrontend(router *gin.Engine) {
	dist, err := fs.Sub(frontend.Dist, "dist")
	if err != nil {
		return
	}

	router.NoRoute(func(c *gin.Context) {
		requestPath := c.Request.URL.Path
		if isReservedBackendPath(requestPath) {
			support.RespondError(c, support.NewError(http.StatusNotFound, "not_found", "route not found"))
			return
		}
		if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodHead {
			support.RespondError(c, support.NewError(http.StatusNotFound, "not_found", "route not found"))
			return
		}

		serveFrontendFile(c, dist, requestPath)
	})
}

func isReservedBackendPath(requestPath string) bool {
	return requestPath == "/api" ||
		strings.HasPrefix(requestPath, "/api/") ||
		requestPath == "/storage" ||
		strings.HasPrefix(requestPath, "/storage/")
}

func serveFrontendFile(c *gin.Context, dist fs.FS, requestPath string) {
	name := strings.TrimPrefix(path.Clean("/"+requestPath), "/")
	if name == "." || name == "" {
		name = "index.html"
	}
	if frontendFileExists(dist, name) {
		c.FileFromFS(name, http.FS(dist))
		return
	}
	c.FileFromFS("index.html", http.FS(dist))
}

func frontendFileExists(dist fs.FS, name string) bool {
	file, err := dist.Open(name)
	if err != nil {
		return false
	}
	defer file.Close()

	stat, err := file.Stat()
	return err == nil && !stat.IsDir()
}
