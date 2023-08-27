package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/custom_error"
	"go.uber.org/zap"
)

type response struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func newResponse(field, message string, err error) response {
	var customError custom_error.CustomError

	if errors.As(err, &customError) {
		resp := response{
			Field:   customError.Field,
			Message: message,
			Error:   customError.Error(),
		}
		return resp
	}

	resp := response{
		Field:   field,
		Message: message,
		Error:   err.Error(),
	}

	return resp
}

func (h *Handler) sentResponse(c *gin.Context, code int, resp response) {
	if resp.Error != "" {
		h.logger.Error(resp.Message, zap.String("errors", resp.Error))
	}
	c.AbortWithStatusJSON(code, resp)
}
