package storage

import "context"

type SegmentStorage interface {
	CreateSegment(ctx context.Context, slug string, percentage uint8) error
}

type Storage interface {
	SegmentStorage
}
