package storage

import "context"

type SegmentStorage interface {
	CreateSegment(ctx context.Context, slug string, percentage uint8) error
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
