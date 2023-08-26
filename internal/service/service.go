package service

//go:generate mockgen -source=service.go -destination=mock/mock.go service

import (
	"context"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/storage"
)

type Segment interface {
	CreateSegment(ctx context.Context, slug string, percentageStr string) error
}

type Services interface {
	Segment
}

type Service struct {
	Segment
}

func NewService(storage storage.Storage) *Service {
	return &Service{
		newSegmentService(storage),
	}
}