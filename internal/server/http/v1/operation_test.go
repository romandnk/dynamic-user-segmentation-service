package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	mock_logger "github.com/romandnk/dynamic-user-segmentation-service/internal/logger/mock"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/service"
	mock_service "github.com/romandnk/dynamic-user-segmentation-service/internal/service/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_CreateCSVReportAndURL(t *testing.T) {
	ctrl := gomock.NewController(t)

	services := mock_service.NewMockServices(ctrl)

	expectedID := uuid.New().String()
	expectedDate := "2023-08"
	expectedUrl := "http://localhost:8080/api/v1/users/report/" + expectedID

	services.EXPECT().CreateCSVReportAndURL(gomock.Any(), expectedDate).Return(expectedUrl, nil)

	handler := NewHandler(services, nil, "")

	r := gin.Default()
	r.POST(url+"/users/report", handler.CreateCSVReportAndURL)

	requestBody := map[string]interface{}{
		"date": expectedDate,
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/users/report", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var responseBody map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	actualURL, ok := responseBody["report_url"]
	require.Equal(t, expectedUrl, actualURL)
	require.True(t, ok)
}

func TestHandler_CreateCSVReportAndURLErrorParsingJSONBody(t *testing.T) {
	ctrl := gomock.NewController(t)

	expectedDate := true

	logger := mock_logger.NewMockLogger(ctrl)

	expectedMessage := "error parsing json body"
	expectedError := "json: cannot unmarshal bool into Go struct field createCSVRepostAndURLBodyRequest.date of type string"

	logger.EXPECT().Error(ErrParsingBody.Error(), zap.String("errors", expectedError))

	handler := NewHandler(nil, logger, "")

	r := gin.Default()
	r.POST(url+"/users/report", handler.CreateCSVReportAndURL)

	requestBody := map[string]interface{}{
		"date": expectedDate,
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/users/report", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusBadRequest, w.Code)

	var responseBody map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	actualMessage, ok := responseBody["message"]
	require.Equal(t, expectedMessage, actualMessage)
	require.True(t, ok)

	actualError, ok := responseBody["error"]
	require.Equal(t, expectedError, actualError)
	require.True(t, ok)
}

func TestHandler_CreateCSVReportAndURLErrorInvalidDate(t *testing.T) {
	ctrl := gomock.NewController(t)

	services := mock_service.NewMockServices(ctrl)
	logger := mock_logger.NewMockLogger(ctrl)

	expectedDate := "2023-08-21"
	expectedError := service.ErrParsingDate
	expectedMessage := "error creating csv report and url"

	services.EXPECT().CreateCSVReportAndURL(gomock.Any(), expectedDate).Return("", expectedError)
	logger.EXPECT().Error(expectedMessage, zap.String("errors", expectedError.Error()))

	handler := NewHandler(services, logger, "")

	r := gin.Default()
	r.POST(url+"/users/report", handler.CreateCSVReportAndURL)

	requestBody := map[string]interface{}{
		"date": expectedDate,
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/users/report", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	var responseBody map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	actualMessage, ok := responseBody["message"]
	require.Equal(t, expectedMessage, actualMessage)
	require.True(t, ok)

	actualError, ok := responseBody["error"]
	require.Equal(t, expectedError.Error(), actualError)
	require.True(t, ok)
}
