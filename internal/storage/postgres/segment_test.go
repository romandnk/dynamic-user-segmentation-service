package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/custom_error"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/models"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestStorage_CreateSegment(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	ctx := context.Background()

	expectedSegment := models.Segment{
		Slug:       "AVITO_TEST",
		Percentage: 10,
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (slug, auto_add_percentage) 
		VALUES ($1, $2)
	`, segmentsTable)

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(expectedSegment.Slug, expectedSegment.Percentage).
		WillReturnResult(pgxmock.NewResult("insert", 1))

	storage := NewStoragePostgres()
	storage.db = mock

	err = storage.CreateSegment(ctx, expectedSegment)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestStorage_CreateSegmentExists(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	ctx := context.Background()

	expectedSegment := models.Segment{
		Slug:       "AVITO_TEST",
		Percentage: 10,
	}
	expectedError := custom_error.CustomError{
		Field:   "slug",
		Message: expectedSegment.Slug + " already exists",
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (slug, auto_add_percentage)
		VALUES ($1, $2)
	`, segmentsTable)

	returnError := &pgconn.PgError{
		Code: "23505",
	}

	mock.ExpectExec(regexp.QuoteMeta(query)).WithArgs(expectedSegment.Slug, expectedSegment.Percentage).
		WillReturnError(returnError)

	storage := NewStoragePostgres()
	storage.db = mock

	err = storage.CreateSegment(ctx, expectedSegment)
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
