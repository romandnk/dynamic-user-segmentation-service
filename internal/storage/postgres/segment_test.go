package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/custom_error"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestStorage_CreateSegment(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	ctx := context.Background()

	expectedSlug := "AVITO_TEST"
	expectedPercentage := uint8(10)

	query := fmt.Sprintf(`
		INSERT INTO %s (slug, auto_add_percentage) 
		VALUES ($1, $2)
	`, segmentsTable)

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(expectedSlug, expectedPercentage).
		WillReturnResult(pgxmock.NewResult("insert", 1))

	storage := NewStoragePostgres()
	storage.db = mock

	err = storage.CreateSegment(ctx, expectedSlug, expectedPercentage)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestStorage_CreateSegmentExists(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	ctx := context.Background()

	expectedSlug := "AVITO_TEST"
	expectedPercentage := uint8(10)
	expectedError := custom_error.CustomError{
		Field:   "slug",
		Message: expectedSlug + " already exists",
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (slug, auto_add_percentage)
		VALUES ($1, $2)
	`, segmentsTable)

	returnError := &pgconn.PgError{
		Code:           "23505",
		Message:        "duplicate key value violates unique constraint \"segments_slug_key\"",
		ConstraintName: "segments_slug_key",
	}

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(expectedSlug, expectedPercentage).WillReturnError(returnError)

	storage := NewStoragePostgres()
	storage.db = mock

	err = storage.CreateSegment(ctx, expectedSlug, expectedPercentage)
	require.ErrorIs(t, err, expectedError)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestStorage_DeleteSegmentOnlyInSegmentsTable(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	ctx := context.Background()

	expectedSlug := "AVITO_TEST"

	queryDeleteFromUserSegments := fmt.Sprintf(`
		DELETE FROM %s
		WHERE segment_slug = $1
		RETURNING user_id
	`, userSegmentsTable)

	queryDeleteFromSegments := fmt.Sprintf(`
		DELETE FROM %s
		WHERE slug = $1
	`, segmentsTable)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryDeleteFromUserSegments)).WithArgs(expectedSlug).
		WillReturnRows(pgxmock.NewRows([]string{}))
	mock.ExpectExec(regexp.QuoteMeta(queryDeleteFromSegments)).WithArgs(expectedSlug).
		WillReturnResult(pgxmock.NewResult("delete", 1))
	mock.ExpectCommit()

	storage := NewStoragePostgres()
	storage.db = mock

	err = storage.DeleteSegment(ctx, expectedSlug)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestStorage_DeleteSegmentInSegmentsTableAndUserSegmentsTable(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	ctx := context.Background()

	expectedSlug := "AVITO_TEST"

	queryDeleteFromUserSegments := fmt.Sprintf(`
		DELETE FROM %s
		WHERE segment_slug = $1
		RETURNING user_id
	`, userSegmentsTable)

	queryDeleteFromSegments := fmt.Sprintf(`
		DELETE FROM %s
		WHERE slug = $1
	`, segmentsTable)

	queryAddOperations := fmt.Sprintf(`
		INSERT INTO %s (user_id, segment_slug, date, operation, auto_add) 
		VALUES ($1, $2, $3, $4, $5)
	`, operationsTable)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryDeleteFromUserSegments)).WithArgs(expectedSlug).
		WillReturnRows(pgxmock.NewRows([]string{"user_id"}).AddRow(1))
	mock.ExpectExec(regexp.QuoteMeta(queryDeleteFromSegments)).WithArgs(expectedSlug).
		WillReturnResult(pgxmock.NewResult("delete", 1))
	mock.ExpectExec(regexp.QuoteMeta(queryAddOperations)).
		WithArgs(1, expectedSlug, pgxmock.AnyArg(), "delete", false).
		WillReturnResult(pgxmock.NewResult("insert", 1))
	mock.ExpectCommit()

	storage := NewStoragePostgres()
	storage.db = mock

	err = storage.DeleteSegment(ctx, expectedSlug)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestStorage_DeleteSegmentNotExist(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	ctx := context.Background()

	expectedSlug := "AVITO_TEST"
	expectedError := custom_error.CustomError{
		Field:   "slug",
		Message: expectedSlug + " doesn't exist",
	}

	queryDeleteFromUserSegments := fmt.Sprintf(`
		DELETE FROM %s
		WHERE segment_slug = $1
		RETURNING user_id
	`, userSegmentsTable)

	queryDeleteFromSegments := fmt.Sprintf(`
		DELETE FROM %s
		WHERE slug = $1
	`, segmentsTable)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryDeleteFromUserSegments)).WithArgs(expectedSlug).
		WillReturnRows(pgxmock.NewRows([]string{}))
	mock.ExpectExec(regexp.QuoteMeta(queryDeleteFromSegments)).WithArgs(expectedSlug).
		WillReturnResult(pgxmock.NewResult("delete", 0))
	mock.ExpectRollback()

	storage := NewStoragePostgres()
	storage.db = mock

	err = storage.DeleteSegment(ctx, expectedSlug)
	require.ErrorIs(t, err, expectedError)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}
