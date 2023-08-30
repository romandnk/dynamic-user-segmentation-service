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

type Operations interface {
	CreateCSVReportAndURL(ctx context.Context, date string) (string, error)
	//GetOperationsByID(ctx context.Context, id string) ([]models.Operation, error)
}

type Services interface {
	Segment
	User
	Operations
}

type Service struct {
	Segment
	User
	Operations
}

func NewService(storage storage.Storage, pathToReports string) *Service {
	return &Service{
		newSegmentService(storage),
		newUserService(storage),
		newOperationService(storage, pathToReports),
	}
}
