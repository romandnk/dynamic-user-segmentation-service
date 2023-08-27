package v1

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/custom_error"
	"net/http"
)

var (
	ErrParsingBody = errors.New("error parsing json body")
)

type createSegmentBodyRequest struct {
	Slug       string `json:"slug"`
	Percentage string `json:"auto_add_percentage"`
}

func (h *Handler) CreateSegment(c *gin.Context) {
	var segmentBody createSegmentBodyRequest

	if err := c.ShouldBindJSON(&segmentBody); err != nil {
		resp := newResponse("", ErrParsingBody.Error(), err)
		h.sentResponse(c, http.StatusBadRequest, resp)
		return
	}

	err := h.services.CreateSegment(c, segmentBody.Slug, segmentBody.Percentage)
	if err != nil {
		message := "error creating segment"
		code := http.StatusInternalServerError
		var customError custom_error.CustomError
		if errors.As(err, &customError) {
			code = http.StatusBadRequest
		}
		resp := newResponse("", message, err)
		h.sentResponse(c, code, resp)
		return
	}

	c.Status(http.StatusCreated)
}

type deleteSegmentBodyRequest struct {
	Slug string `json:"slug"`
}

func (h *Handler) DeleteSegment(c *gin.Context) {
	var segmentBody deleteSegmentBodyRequest

	if err := c.ShouldBindJSON(&segmentBody); err != nil {
		resp := newResponse("", ErrParsingBody.Error(), err)
		h.sentResponse(c, http.StatusBadRequest, resp)
		return
	}

	err := h.services.DeleteSegment(c, segmentBody.Slug)
	if err != nil {
		message := "error deleting segment"
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
