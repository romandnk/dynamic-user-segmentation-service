package service

//go:generate mockgen -source=service.go -destination=mock/mock.go service

import (
	"context"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/storage"
)

type Segment interface {
	CreateSegment(ctx context.Context, slug string, percentageStr string) error
	DeleteSegment(ctx context.Context, slug string) error
}

type User interface {
	UpdateUserSegments(ctx context.Context, segmentsToAdd, segmentsToDelete []string, userID int) error
	GetActiveSegments(ctx context.Context, userID int) ([]string, error)
	AutoAddSegments(ctx context.Context) error
}

type Services interface {
	Segment
	User
}

type Service struct {
	Segment
	User
}

func NewService(storage storage.Storage) *Service {
	return &Service{
		newSegmentService(storage),
		newUserService(storage),
	}
}
