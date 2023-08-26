package storage

import "context"

type SegmentStorage interface {
	CreateSegment(ctx context.Context, slug string, percentage uint8) error
	DeleteSegment(ctx context.Context, slug string) error
}

type Storage interface {
	SegmentStorage
}
