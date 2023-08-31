package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	_ "github.com/romandnk/dynamic-user-segmentation-service/docs"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/custom_error"
	"net/http"
)

type addAndDeleteUserSegmentsBodyRequest struct {
	SegmentsToAdd    []string `json:"segments_to_add"`
	SegmentsToDelete []string `json:"segments_to_delete"`
	UserID           int      `json:"user_id"`
}

// UpdateUserSegments godoc
// @Summary Add and delete user segments by his id
// @Tags user
// @Accept json
// @Param input body addAndDeleteUserSegmentsBodyRequest true "user segments to add and delete and his user id"
// @Success 200
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /users [post]
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

type getActiveUserSegmentsBodyResponse struct {
	Segments []string `json:"segments"`
}

// GetActiveUserSegments godoc
// @Summary Get active user segments
// @Tags user
// @Accept json
// @Param input body getActiveUserSegmentsBodyRequest true "user id to get his segments"
// @Success 200 {object} getActiveUserSegmentsBodyResponse
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /users/active_segments [post]
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

	c.JSON(http.StatusOK, getActiveUserSegmentsBodyResponse{
		Segments: segments,
	})
}
