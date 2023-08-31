package postgres

import (
	"context"
	"fmt"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/models"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

func TestStorage_GetOperations(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	ctx := context.Background()

	expectedDate := time.Date(2023, 8, 1, 0, 0, 0, 0, time.Local)
	expectedOperations := []models.Operation{
		{
			UserID:      1,
			SegmentSlug: "TEST",
			Date:        expectedDate,
			Action:      "add",
		},
	}

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

	columns := []string{"user_id", "segment_slug", "date", "action"}
	rows := pgxmock.NewRows(columns).AddRow(1, "TEST", expectedDate, "add")

	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WithArgs(expectedDate, expectedDate.AddDate(0, 1, 0).Add(-time.Nanosecond)).
		WillReturnRows(rows)

	storage := NewStoragePostgres()
	storage.db = mock

	operations, err := storage.GetOperations(ctx, expectedDate)
	require.NoError(t, err)
	require.ElementsMatch(t, operations, expectedOperations)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}
