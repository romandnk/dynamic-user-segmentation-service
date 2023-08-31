package v1

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/romandnk/dynamic-user-segmentation-service/docs"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/custom_error"
	"net/http"
)

type createCSVRepostAndURLBodyRequest struct {
	Date string `json:"date"`
}

type createCSVRepostAndURLBodyResponse struct {
	URL string `json:"report_url"`
}

// CreateCSVReportAndURL godoc
// @Summary Create a CSV file locally and return url to download a file
// @Tags operation
// @Accept json
// @Param input body createCSVRepostAndURLBodyRequest true "date format year-month"
// @Success 200 {object} createCSVRepostAndURLBodyResponse
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /users/report [post]
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

	c.JSON(http.StatusOK, createCSVRepostAndURLBodyResponse{
		URL: url,
	})
}

// GetReportByID godoc
// @Summary Get report CSV file to download
// @Tags operation
// @Param input path string true "report id"
// @Success 200
// @Failure 400 {object} response
// @Failure 500 {object} response
// @Router /users/report/{id} [get]
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
