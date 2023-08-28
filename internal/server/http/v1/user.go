package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/custom_error"
	"net/http"
)

type addAndDeleteUserSegmentsBodyRequest struct {
	SegmentsToAdd    []string `json:"segments_to_add"`
	SegmentsToDelete []string `json:"segments_to_delete"`
	UserID           int      `json:"user_id"`
}

func (h *Handler) AddAndDeleteUserSegments(c *gin.Context) {
	var addAndDeleteUserSegmentsBody addAndDeleteUserSegmentsBodyRequest

	if err := c.ShouldBindJSON(&addAndDeleteUserSegmentsBody); err != nil {
		resp := newResponse("", ErrParsingBody.Error(), err)
		h.sentResponse(c, http.StatusBadRequest, resp)
		return
	}

	err := h.services.UpdateUserSegments(c,
		addAndDeleteUserSegmentsBody.SegmentsToAdd,
		addAndDeleteUserSegmentsBody.SegmentsToDelete,
		addAndDeleteUserSegmentsBody.UserID,
	)
	if err != nil {
		message := "error updating user segments"
		code := http.StatusInternalServerError
		var customError custom_error.CustomError
		if errors.As(err, &customError) {
			code = http.StatusBadRequest
		}
		resp := newResponse("", message, err)
		h.sentResponse(c, code, resp)
		return
	}

	c.Status(http.StatusOK)
}
