package service

import (
	"github.com/romandnk/dynamic-user-segmentation-service/internal/custom_error"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidatePercentage(t *testing.T) {
	testCases := []struct {
		name           string
		input          string
		expectedOutput uint8
		expectedError  error
	}{
		{
			name:           "valid percentage",
			input:          "50%",
			expectedOutput: 50,
			expectedError:  nil,
		},
		{
			name:           "invalid percentage zero",
			input:          "0%",
			expectedOutput: 0,
			expectedError: custom_error.CustomError{
				Field:   "percentage",
				Message: ErrInvalidPercentageZero.Error(),
			},
		},
		{
			name:           "input without percent (%)",
			input:          "50",
			expectedOutput: 0,
			expectedError: custom_error.CustomError{
				Field:   "percentage",
				Message: ErrInvalidPercentageFormat.Error(),
			},
		},
		{
			name:           "percent more than 100",
			input:          "200%",
			expectedOutput: 0,
			expectedError: custom_error.CustomError{
				Field:   "percentage",
				Message: ErrInvalidPercentageTooBig.Error(),
			},
		},
		{
			name:           "percent less than 0",
			input:          "-2%",
			expectedOutput: 0,
			expectedError: custom_error.CustomError{
				Field:   "percentage",
				Message: ErrInvalidPercentageFormat.Error(),
			},
		},
		{
			name:           "empty input",
			input:          "",
			expectedOutput: 0,
			expectedError:  nil,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			actualOutput, actualError := validatePercentage(tc.input)
			require.Equal(t, tc.expectedOutput, actualOutput)
			require.ErrorIs(t, actualError, tc.expectedError)
		})
	}
}
