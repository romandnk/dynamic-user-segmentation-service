package service

import (
	"context"
	"errors"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/custom_error"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/models"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/storage"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrEmptySlug                 = errors.New("empty slug")
	ErrInvalidSlugRepresentation = errors.New("slug can only contain uppercase letters")
	ErrInvalidPercentageZero     = errors.New("percentage cannot be zero")
	ErrInvalidPercentageFormat   = errors.New("invalid percentage format (e.g. 100%, 99%, 1%)")
	ErrInvalidPercentageTooBig   = errors.New("percentage cannot be more than 100")
)

var percentageValidFormat = regexp.MustCompile(`^\d+%$`)

type segmentService struct {
	segment storage.SegmentStorage
}

func newSegmentService(segment storage.SegmentStorage) *segmentService {
	return &segmentService{segment: segment}
}

func (s *segmentService) CreateSegment(ctx context.Context, slug string, percentageStr string) error {
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

	segment := models.Segment{
		Slug:       slug,
		Percentage: percentage,
	}

	return s.segment.CreateSegment(ctx, segment)
}

func validatePercentage(percentageStr string) (int, error) {
	if percentageStr == "" {
		return 0, nil
	}

	if !percentageValidFormat.MatchString(percentageStr) {
		return 0, custom_error.CustomError{
			Field:   "percentage",
			Message: ErrInvalidPercentageFormat.Error(),
		}
	}

	percentage, err := strconv.ParseInt(percentageStr[:len(percentageStr)-1], 10, 64)
	if err != nil {
		return 0, custom_error.CustomError{
			Field:   "percentage",
			Message: err.Error(),
		}
	}

	if percentage > 100 {
		return 0, custom_error.CustomError{
			Field:   "percentage",
			Message: ErrInvalidPercentageTooBig.Error(),
		}
	}

	if percentage == 0 {
		return 0, custom_error.CustomError{
			Field:   "percentage",
			Message: ErrInvalidPercentageZero.Error(),
		}
	}

	return int(percentage), nil
}

func (s *segmentService) DeleteSegment(ctx context.Context, slug string) error {
	slug = strings.TrimSpace(slug)

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

	return s.segment.DeleteSegment(ctx, slug)
}
