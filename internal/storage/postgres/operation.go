package postgres

import (
	"context"
	"fmt"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/models"
	"time"
)

func (s *Storage) GetOperations(ctx context.Context, date time.Time) ([]models.Operation, error) {
	query := fmt.Sprintf(`
		SELECT
    		user_id,
    		segment_slug,
    		date,
    		action
		FROM %s
		WHERE date >= $1 AND date <= $2
		ORDER BY user_id
	`, operationsTable)

	rows, err := s.db.Query(ctx, query, date, date.AddDate(0, 1, 0).Add(-time.Nanosecond))
	if err != nil {
		return nil, fmt.Errorf("OperationRepo.GetOperations - s.db.Query: %w", err)
	}
	defer rows.Close()

	var operations []models.Operation
	for rows.Next() {
		var operation models.Operation

		err = rows.Scan(&operation.UserID, &operation.SegmentSlug, &operation.Date, &operation.Action)
		if err != nil {
			return nil, fmt.Errorf("OperationRepo.GetOperations - rows.Scan: %w", err)
		}

		operations = append(operations, operation)
	}

	return operations, nil
}
