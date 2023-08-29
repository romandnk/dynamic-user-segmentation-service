package storage

import (
	"context"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/models"
)

type SegmentStorage interface {
	CreateSegment(ctx context.Context, segment models.Segment) error
	DeleteSegment(ctx context.Context, slug string) error
}

type UserStorage interface {
	UpdateUserSegments(ctx context.Context, segmentsToAdd, segmentsToDelete []string, userID int, random uint8) error
	GetActiveSegments(ctx context.Context, userID int) ([]string, error)
}

type Storage interface {
	SegmentStorage
	UserStorage
}
