package storage

import "context"

type SegmentStorage interface {
	CreateSegment(ctx context.Context, slug string, percentage uint8) error
	DeleteSegment(ctx context.Context, slug string) error
}

type UserStorage interface {
	UpdateUserSegments(ctx context.Context, segmentsToAdd, segmentsToDelete []string, id int) error
}

type Storage interface {
	SegmentStorage
	UserStorage
}
