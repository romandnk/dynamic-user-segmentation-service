package service

import (
	"context"
	"errors"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/custom_error"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/storage"
	"math/rand"
	"strings"
	"time"
)

var (
	ErrBothEmptySegments            = errors.New("segments to add and segments to delete cannot both be empty")
	ErrInvalidSegmentRepresentation = errors.New("segment can only contain uppercase letters")
	ErrInvalidUserID                = errors.New("user id can be only positive number")
)

type userService struct {
	user storage.UserStorage
}

func newUserService(user storage.UserStorage) *userService {
	return &userService{user: user}
}

func (u *userService) UpdateUserSegments(ctx context.Context, segmentsToAdd, segmentsToDelete []string, userID int) error {
	if userID <= 0 {
		return custom_error.CustomError{
			Field:   "user_id",
			Message: ErrInvalidUserID.Error(),
		}
	}

	if len(segmentsToAdd) == 0 && len(segmentsToDelete) == 0 {
		return custom_error.CustomError{
			Field:   "segments",
			Message: ErrBothEmptySegments.Error(),
		}
	}

	for _, segment := range segmentsToAdd {
		if strings.ToUpper(segment) != segment {
			return custom_error.CustomError{
				Field:   "segment to add",
				Message: ErrInvalidSegmentRepresentation.Error(),
			}
		}
	}

	for _, segment := range segmentsToDelete {
		if strings.ToUpper(segment) != segment {
			return custom_error.CustomError{
				Field:   "segment to delete",
				Message: ErrInvalidSegmentRepresentation.Error(),
			}
		}
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	random := 1 + r.Intn(100)

	return u.user.UpdateUserSegments(ctx, segmentsToAdd, segmentsToDelete, userID, uint8(random))
}

func (u *userService) GetActiveSegments(ctx context.Context, userID int) ([]string, error) {
	if userID <= 0 {
		return nil, custom_error.CustomError{
			Field:   "user_id",
			Message: ErrInvalidUserID.Error(),
		}
	}

	return u.user.GetActiveSegments(ctx, userID)
}
