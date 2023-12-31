package storage

import (
	"context"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/models"
	"time"
)

type SegmentStorage interface {
	CreateSegment(ctx context.Context, segment models.Segment) error
	DeleteSegment(ctx context.Context, slug string) error
}

type UserStorage interface {
	UpdateUserSegments(ctx context.Context, segmentsToAdd, segmentsToDelete []string, userID int) error
	GetActiveSegments(ctx context.Context, userID int) ([]string, error)
	AutoAddUserSegments(ctx context.Context) error
}

type OperationStorage interface {
	GetOperations(ctx context.Context, date time.Time) ([]models.Operation, error)
}

type Storage interface {
	SegmentStorage
	UserStorage
	OperationStorage
}
