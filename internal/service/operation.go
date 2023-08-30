package service

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/models"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/storage"
	"github.com/spf13/viper"
	"net"
	"net/url"
	"os"
	"strconv"
	"time"
)

var (
	ErrParsingDate = errors.New("invalid date (year-month, e.g. 2023-08)")
)

type operationService struct {
	operation     storage.OperationStorage
	pathToReports string
}

func newOperationService(operation storage.OperationStorage, pathToReports string) *operationService {
	return &operationService{
		operation:     operation,
		pathToReports: pathToReports,
	}
}

//func (o *operationService) GetOperationsByID(ctx context.Context, id string) ([]models.Operation, error) {
//	return o.operation.GetOperations(ctx, time.Time{})
//}

func (o *operationService) CreateCSVReportAndURL(ctx context.Context, date string) (string, error) {
	layout := "2006-01"
	parsedTime, err := time.Parse(layout, date)
	if err != nil {
		return "", ErrParsingDate
	}

	operations, err := o.operation.GetOperations(ctx, parsedTime)
	if err != nil {
		return "", err
	}

	id := uuid.New().String()

	err = createCSVFile(o.pathToReports, operations, id)
	if err != nil {
		return "", err
	}

	host, err := os.Hostname()
	if err != nil {
		return "", err
	}

	IPs, err := net.LookupIP(host)
	if err != nil {
		return "", err
	}

	port := viper.GetString("server.port")

	u := url.URL{
		Scheme: "http",
		Host:   net.JoinHostPort(IPs[0].String(), port),
		Path:   fmt.Sprintf("/api/v1/users/report/%s", id),
	}

	return u.String(), nil
}

func createCSVFile(path string, operations []models.Operation, id string) error {
	f, err := os.Create(path + id + ".csv")
	if err != nil {
		return fmt.Errorf("error creating file with id %s: %w", id, err)
	}
	defer func() {
		_ = f.Close()
	}()

	w := csv.NewWriter(f)
	defer w.Flush()

	columns := []string{"user id", "segment_slug", "action", "date"}
	err = w.Write(columns)
	if err != nil {
		return fmt.Errorf("error writing column to file with id %s: %w", id, err)
	}

	for _, operation := range operations {
		userID := strconv.Itoa(operation.UserID)
		segmentSlug := operation.SegmentSlug
		action := operation.Action
		date := operation.Date.Format(time.DateTime) // a human-readable format
		row := []string{userID, segmentSlug, action, date}

		if err := w.Write(row); err != nil {
			return fmt.Errorf("error writing file with id %s: %w", id, err)
		}
	}

	return nil
}
