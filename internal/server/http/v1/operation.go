package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/custom_error"
	"net/http"
)

type createCSVRepostAndURLBodyRequest struct {
	Date string `json:"date"`
}

func (h *Handler) CreateCSVReportAndURL(c *gin.Context) {
	var createCSVRepostAndURLBody createCSVRepostAndURLBodyRequest

	if err := c.ShouldBindJSON(&createCSVRepostAndURLBody); err != nil {
		resp := newResponse("", ErrParsingBody.Error(), err)
		h.sentResponse(c, http.StatusBadRequest, resp)
		return
	}

	url, err := h.services.CreateCSVReportAndURL(c, createCSVRepostAndURLBody.Date)
	if err != nil {
		message := "error creating csv report and url"
		code := http.StatusInternalServerError
		var customError custom_error.CustomError
		if errors.As(err, &customError) {
			code = http.StatusBadRequest
		}
		resp := newResponse("", message, err)
		h.sentResponse(c, code, resp)
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"report_url": url,
	})
}

func (h *Handler) GetReportByID(c *gin.Context) {
	id := c.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		resp := newResponse("", ErrParsingBody.Error(), err)
		h.sentResponse(c, http.StatusBadRequest, resp)
		return
	}

	fileName := parsedID.String() + ".csv"
	filePath := h.pathToReports + fileName

	fmt.Println(filePath)

	c.FileAttachment(filePath, fileName)
}
