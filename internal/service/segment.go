package service

import (
	"context"
	"errors"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/custom_error"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/storage"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrEmptySlug                 = errors.New("empty slug")
	ErrInvalidSlugRepresentation = errors.New("slug can only contain uppercase letters")
	ErrInvalidPercentageZero     = errors.New("percentage cannot be zero")
	ErrInvalidPercentageFormat   = errors.New("invalid slug format (e.g. 100%, 99%, 1%)")
)

var percentageValidFormat = regexp.MustCompile(`^(100|\d{1,2})%$`)

type SegmentService struct {
	segment storage.SegmentStorage
}

func NewSegmentService(segment storage.SegmentStorage) *SegmentService {
	return &SegmentService{segment: segment}
}

func (s *SegmentService) CreateSegment(ctx context.Context, slug string, percentageStr string) error {
	slug = strings.TrimSpace(slug)
	percentageStr = strings.TrimSpace(percentageStr)

	if slug == "" {
		return custom_error.CustomError{
			Field:   "slug",
			Message: ErrEmptySlug.Error(),
		}
	}

	if strings.ToUpper(slug) != slug {
		return custom_error.CustomError{
			Field:   "slug",
			Message: ErrInvalidSlugRepresentation.Error(),
		}
	}

	percentage, err := validatePercentage(percentageStr)
	if err != nil {
		return err
	}

	return s.segment.CreateSegment(ctx, slug, percentage)
}

func validatePercentage(percentageStr string) (uint8, error) {
	if percentageStr == "" {
		return 0, nil
	}

	if !percentageValidFormat.MatchString(percentageStr) {
		return 0, custom_error.CustomError{
			Field:   "percentage",
			Message: ErrInvalidPercentageFormat.Error(),
		}
	}

	percentage, err := strconv.ParseUint(percentageStr[:len(percentageStr)-1], 10, 8)
	if err != nil {
		return 0, custom_error.CustomError{
			Field:   "percentage",
			Message: err.Error(),
		}
	}

	if percentage == 0 {
		return 0, custom_error.CustomError{
			Field:   "percentage",
			Message: ErrInvalidPercentageZero.Error(),
		}
	}

	return uint8(percentage), nil
}
