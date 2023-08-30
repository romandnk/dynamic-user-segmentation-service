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

func (h *Handler) UpdateUserSegments(c *gin.Context) {
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

type getActiveUserSegmentsBodyRequest struct {
	UserID int `json:"user_id"`
}

func (h *Handler) GetActiveUserSegments(c *gin.Context) {
	var getActiveUserSegmentsBody getActiveUserSegmentsBodyRequest

	if err := c.ShouldBindJSON(&getActiveUserSegmentsBody); err != nil {
		resp := newResponse("", ErrParsingBody.Error(), err)
		h.sentResponse(c, http.StatusBadRequest, resp)
		return
	}

	segments, err := h.services.GetActiveSegments(c, getActiveUserSegmentsBody.UserID)
	if err != nil {
		message := "error getting user segments"
		code := http.StatusInternalServerError
		var customError custom_error.CustomError
		if errors.As(err, &customError) {
			code = http.StatusBadRequest
		}
		resp := newResponse("", message, err)
		h.sentResponse(c, code, resp)
		return
	}

	if len(segments) == 0 {
		c.JSON(http.StatusOK, map[string]string{
			"segments": "no segments",
		})
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"segments": segments,
	})
}
