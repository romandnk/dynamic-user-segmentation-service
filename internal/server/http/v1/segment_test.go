package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/custom_error"
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

const url = "/api/v1"

func TestHandler_CreateSegment(t *testing.T) {
	ctrl := gomock.NewController(t)

	services := mock_service.NewMockServices(ctrl)

	expectedSlug := "AVITO_TEST"
	expectedAutoAddPercentage := "10%"

	services.EXPECT().CreateSegment(gomock.Any(), expectedSlug, expectedAutoAddPercentage).Return(nil)

	handler := NewHandler(services, nil, "")

	r := gin.Default()
	r.POST(url+"/segments", handler.CreateSegment)

	requestBody := map[string]interface{}{
		"slug":                "AVITO_TEST",
		"auto_add_percentage": "10%",
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/segments", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	require.Equal(t, []byte(nil), w.Body.Bytes())
}

func TestHandler_CreateSegmentErrorParsingJSONBody(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := mock_logger.NewMockLogger(ctrl)

	expectedMessage := "error parsing json body"
	expectedError := "json: cannot unmarshal bool into Go struct field createSegmentBodyRequest.slug of type string"

	logger.EXPECT().Error(ErrParsingBody.Error(), zap.String("errors", expectedError))

	handler := NewHandler(nil, logger, "")

	r := gin.Default()
	r.POST(url+"/segments", handler.CreateSegment)

	requestBody := map[string]interface{}{
		"slug": true,
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/segments", bytes.NewBuffer(jsonBody))
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

func TestHandler_CreateSegmentError(t *testing.T) {
	expectedMessage := "error creating segment"

	testCases := []struct {
		name            string
		inputSlug       string
		inputPercentage string
		expectedError   error
		expectedField   string
		expectedCode    int
	}{
		{
			name:            "slug in lowercase",
			inputSlug:       "avito_small",
			inputPercentage: "",
			expectedError: custom_error.CustomError{
				Field:   "slug",
				Message: service.ErrInvalidSlugRepresentation.Error(),
			},
			expectedField: "slug",
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:            "slug is empty",
			inputSlug:       "",
			inputPercentage: "",
			expectedError: custom_error.CustomError{
				Field:   "slug",
				Message: service.ErrEmptySlug.Error(),
			},
			expectedField: "slug",
			expectedCode:  http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			services := mock_service.NewMockServices(ctrl)
			logger := mock_logger.NewMockLogger(ctrl)

			logger.EXPECT().Error(expectedMessage, zap.String("errors", tc.expectedError.Error()))
			services.EXPECT().CreateSegment(gomock.Any(), tc.inputSlug, tc.inputPercentage).Return(tc.expectedError)

			handler := NewHandler(services, logger, "")

			r := gin.Default()
			r.POST(url+"/segments", handler.CreateSegment)

			requestBody := map[string]interface{}{
				"slug":                tc.inputSlug,
				"auto_add_percentage": tc.inputPercentage,
			}

			jsonBody, err := json.Marshal(requestBody)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/segments", bytes.NewBuffer(jsonBody))

			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			require.Equal(t, tc.expectedCode, w.Code)

			var responseBody map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			require.NoError(t, err)

			actualField, ok := responseBody["field"]
			require.Equal(t, tc.expectedField, actualField)
			require.True(t, ok)

			actualMessage, ok := responseBody["message"]
			require.Equal(t, expectedMessage, actualMessage)
			require.True(t, ok)

			actualError, ok := responseBody["error"]
			require.Equal(t, tc.expectedError.Error(), actualError)
			require.True(t, ok)
		})
	}
}

func TestHandler_DeleteSegment(t *testing.T) {
	ctrl := gomock.NewController(t)

	services := mock_service.NewMockServices(ctrl)

	expectedSlug := "AVITO_TEST"

	services.EXPECT().DeleteSegment(gomock.Any(), expectedSlug).Return(nil)

	handler := NewHandler(services, nil, "")

	r := gin.Default()
	r.DELETE(url+"/segments", handler.DeleteSegment)

	requestBody := map[string]interface{}{
		"slug": "AVITO_TEST",
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url+"/segments", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	require.Equal(t, []byte(nil), w.Body.Bytes())
}

func TestHandler_DeleteSegmentErrorParsingJsonBody(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := mock_logger.NewMockLogger(ctrl)

	expectedMessage := "error parsing json body"
	expectedError := "json: cannot unmarshal bool into Go struct field deleteSegmentBodyRequest.slug of type string"

	logger.EXPECT().Error(ErrParsingBody.Error(), zap.String("errors", expectedError))

	handler := NewHandler(nil, logger, "")

	r := gin.Default()
	r.DELETE(url+"/segments", handler.DeleteSegment)

	requestBody := map[string]interface{}{
		"slug": true,
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url+"/segments", bytes.NewBuffer(jsonBody))
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

func TestHandler_DeleteSegmentError(t *testing.T) {
	expectedMessage := "error deleting segment"

	testCases := []struct {
		name          string
		inputSlug     string
		expectedError error
		expectedField string
		expectedCode  int
	}{
		{
			name:      "slug in lowercase",
			inputSlug: "avito_small",
			expectedError: custom_error.CustomError{
				Field:   "slug",
				Message: service.ErrInvalidSlugRepresentation.Error(),
			},
			expectedField: "slug",
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:      "slug is empty",
			inputSlug: "",
			expectedError: custom_error.CustomError{
				Field:   "slug",
				Message: service.ErrEmptySlug.Error(),
			},
			expectedField: "slug",
			expectedCode:  http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			services := mock_service.NewMockServices(ctrl)
			logger := mock_logger.NewMockLogger(ctrl)

			logger.EXPECT().Error(expectedMessage, zap.String("errors", tc.expectedError.Error()))
			services.EXPECT().DeleteSegment(gomock.Any(), tc.inputSlug).Return(tc.expectedError)

			handler := NewHandler(services, logger, "")

			r := gin.Default()
			r.DELETE(url+"/segments", handler.DeleteSegment)

			requestBody := map[string]interface{}{
				"slug": tc.inputSlug,
			}

			jsonBody, err := json.Marshal(requestBody)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url+"/segments", bytes.NewBuffer(jsonBody))

			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			require.Equal(t, tc.expectedCode, w.Code)

			var responseBody map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &responseBody)
			require.NoError(t, err)

			actualField, ok := responseBody["field"]
			require.Equal(t, tc.expectedField, actualField)
			require.True(t, ok)

			actualMessage, ok := responseBody["message"]
			require.Equal(t, expectedMessage, actualMessage)
			require.True(t, ok)

			actualError, ok := responseBody["error"]
			require.Equal(t, tc.expectedError.Error(), actualError)
			require.True(t, ok)
		})
	}
}
