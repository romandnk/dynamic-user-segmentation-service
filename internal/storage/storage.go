package storage

import "context"

type SegmentStorage interface {
	CreateSegment(ctx context.Context, slug string, percentage float32) (int, error)
}

type Storage interface {
	SegmentStorage
}
