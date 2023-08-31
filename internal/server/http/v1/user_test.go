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

func TestHandler_UpdateUserSegments(t *testing.T) {
	ctrl := gomock.NewController(t)

	services := mock_service.NewMockServices(ctrl)

	expectedSegmentsToAdd := []string{"AVITO_TEST1", "AVITO_TEST2"}
	expectedSegmentsToDelete := []string{"AVITO_TEST3"}
	expectedUserID := 1

	services.EXPECT().UpdateUserSegments(gomock.Any(), expectedSegmentsToAdd, expectedSegmentsToDelete, expectedUserID).
		Return(nil)

	handler := NewHandler(services, nil, "")

	r := gin.Default()
	r.POST(url+"/users", handler.UpdateUserSegments)

	requestBody := map[string]interface{}{
		"segments_to_add":    expectedSegmentsToAdd,
		"segments_to_delete": expectedSegmentsToDelete,
		"user_id":            expectedUserID,
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/users", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	require.Equal(t, []byte(nil), w.Body.Bytes())
}

func TestHandler_UpdateUserSegmentsErrorParsingJSONBody(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := mock_logger.NewMockLogger(ctrl)

	expectedMessage := "error parsing json body"
	expectedError := "json: cannot unmarshal bool into Go struct field addAndDeleteUserSegmentsBodyRequest.segments_to_add of type []string"

	logger.EXPECT().Error(ErrParsingBody.Error(), zap.String("errors", expectedError))

	handler := NewHandler(nil, logger, "")

	r := gin.Default()
	r.POST(url+"/users", handler.UpdateUserSegments)

	requestBody := map[string]interface{}{
		"segments_to_add":    true,
		"segments_to_delete": []string{},
		"user_id":            1,
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/users", bytes.NewBuffer(jsonBody))
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

func TestHandler_UpdateUserSegmentsError(t *testing.T) {
	expectedMessage := "error updating user segments"

	testCases := []struct {
		name                  string
		inputSegmentsToAdd    []string
		inputSegmentsToDelete []string
		inputUserID           int
		expectedError         error
		expectedField         string
		expectedCode          int
	}{
		{
			name:                  "empty both add and delete segments",
			inputSegmentsToAdd:    []string{},
			inputSegmentsToDelete: []string{},
			inputUserID:           1,
			expectedError: custom_error.CustomError{
				Field:   "segments",
				Message: service.ErrInvalidSegmentRepresentation.Error(),
			},
			expectedField: "segments",
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:                  "user id less or equal zero",
			inputSegmentsToAdd:    []string{"TEST"},
			inputSegmentsToDelete: []string{},
			inputUserID:           0,
			expectedError: custom_error.CustomError{
				Field:   "user_id",
				Message: service.ErrInvalidUserID.Error(),
			},
			expectedField: "user_id",
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:                  "segments to add have segment in lowercase",
			inputSegmentsToAdd:    []string{"test"},
			inputSegmentsToDelete: []string{},
			inputUserID:           1,
			expectedError: custom_error.CustomError{
				Field:   "segment to add",
				Message: service.ErrInvalidSegmentRepresentation.Error(),
			},
			expectedField: "segment to add",
			expectedCode:  http.StatusBadRequest,
		},
		{
			name:                  "segments to delete have segment in lowercase",
			inputSegmentsToAdd:    []string{},
			inputSegmentsToDelete: []string{"test"},
			inputUserID:           1,
			expectedError: custom_error.CustomError{
				Field:   "segment to delete",
				Message: service.ErrInvalidSegmentRepresentation.Error(),
			},
			expectedField: "segment to delete",
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
			services.EXPECT().
				UpdateUserSegments(gomock.Any(), tc.inputSegmentsToAdd, tc.inputSegmentsToDelete, tc.inputUserID).
				Return(tc.expectedError)

			handler := NewHandler(services, logger, "")

			r := gin.Default()
			r.POST(url+"/users", handler.UpdateUserSegments)

			requestBody := map[string]interface{}{
				"segments_to_add":    tc.inputSegmentsToAdd,
				"segments_to_delete": tc.inputSegmentsToDelete,
				"user_id":            tc.inputUserID,
			}

			jsonBody, err := json.Marshal(requestBody)
			require.NoError(t, err)

			w := httptest.NewRecorder()

			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/users", bytes.NewBuffer(jsonBody))

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

func TestHandler_GetActiveUserSegments(t *testing.T) {
	ctrl := gomock.NewController(t)

	services := mock_service.NewMockServices(ctrl)

	expectedUserID := 1
	expectedUserSegments := []string{"TEST1", "TEST2"}

	services.EXPECT().GetActiveSegments(gomock.Any(), expectedUserID).Return(expectedUserSegments, nil)

	handler := NewHandler(services, nil, "")

	r := gin.Default()
	r.POST(url+"/users/active_segments", handler.GetActiveUserSegments)

	requestBody := map[string]interface{}{
		"user_id": expectedUserID,
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/users/active_segments", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var responseBody map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	segments, ok := responseBody["segments"]
	require.True(t, ok)
	require.ElementsMatch(t, segments, expectedUserSegments)
}

func TestHandler_GetActiveUserSegmentsEmptySegments(t *testing.T) {
	ctrl := gomock.NewController(t)

	services := mock_service.NewMockServices(ctrl)

	expectedUserID := 1
	expectedUserSegments := []string{}

	services.EXPECT().GetActiveSegments(gomock.Any(), expectedUserID).Return(expectedUserSegments, nil)

	handler := NewHandler(services, nil, "")

	r := gin.Default()
	r.POST(url+"/users/active_segments", handler.GetActiveUserSegments)

	requestBody := map[string]interface{}{
		"user_id": expectedUserID,
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/users/active_segments", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var responseBody map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &responseBody)
	require.NoError(t, err)

	segments, ok := responseBody["segments"]
	require.True(t, ok)
	require.Equal(t, []interface{}{}, segments)
}

func TestHandler_GetActiveUserSegmentsErrorParsingJSONBody(t *testing.T) {
	ctrl := gomock.NewController(t)

	expectedUserID := "text"

	logger := mock_logger.NewMockLogger(ctrl)

	expectedMessage := "error parsing json body"
	expectedError := "json: cannot unmarshal string into Go struct field getActiveUserSegmentsBodyRequest.user_id of type int"

	logger.EXPECT().Error(ErrParsingBody.Error(), zap.String("errors", expectedError))

	handler := NewHandler(nil, logger, "")

	r := gin.Default()
	r.POST(url+"/users/active_segments", handler.GetActiveUserSegments)

	requestBody := map[string]interface{}{
		"user_id": expectedUserID,
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/users/active_segments", bytes.NewBuffer(jsonBody))
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

func TestHandler_GetActiveUserSegmentsErrorInvalidUserID(t *testing.T) {
	ctrl := gomock.NewController(t)

	services := mock_service.NewMockServices(ctrl)
	logger := mock_logger.NewMockLogger(ctrl)

	expectedUserID := 0
	expectedUserSegments := []string{}
	expectedError := custom_error.CustomError{
		Field:   "user_id",
		Message: service.ErrInvalidUserID.Error(),
	}
	expectedMessage := "error getting user segments"

	services.EXPECT().GetActiveSegments(gomock.Any(), expectedUserID).Return(expectedUserSegments, expectedError)
	logger.EXPECT().Error(expectedMessage, zap.String("errors", expectedError.Error()))

	handler := NewHandler(services, logger, "")

	r := gin.Default()
	r.POST(url+"/users/active_segments", handler.GetActiveUserSegments)

	requestBody := map[string]interface{}{
		"user_id": expectedUserID,
	}

	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	w := httptest.NewRecorder()

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url+"/users/active_segments", bytes.NewBuffer(jsonBody))
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
	require.Equal(t, expectedError.Error(), actualError)
	require.True(t, ok)
}
