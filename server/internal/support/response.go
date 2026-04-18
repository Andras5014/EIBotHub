package support

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Envelope struct {
	Success bool       `json:"success"`
	Message string     `json:"message,omitempty"`
	Data    any        `json:"data,omitempty"`
	Error   *ErrorBody `json:"error,omitempty"`
}

type AppError struct {
	Status  int
	Code    string
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewError(status int, code, message string) *AppError {
	return &AppError{Status: status, Code: code, Message: message}
}

func RespondOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, Envelope{Success: true, Data: data})
}

func RespondCreated(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, Envelope{Success: true, Data: data})
}

func RespondMessage(c *gin.Context, message string, data any) {
	c.JSON(http.StatusOK, Envelope{Success: true, Message: message, Data: data})
}

func RespondError(c *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		c.JSON(appErr.Status, Envelope{
			Success: false,
			Error:   &ErrorBody{Code: appErr.Code, Message: appErr.Message},
		})
		return
	}

	c.JSON(http.StatusInternalServerError, Envelope{
		Success: false,
		Error:   &ErrorBody{Code: "internal_error", Message: "internal server error"},
	})
}

func Paginated(items any, page, pageSize int, total int64) map[string]any {
	return map[string]any{
		"items":     items,
		"page":      page,
		"page_size": pageSize,
		"total":     total,
	}
}
