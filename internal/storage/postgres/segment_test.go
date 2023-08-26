package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/custom_error"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestStorage_CreateSegmentNotExist(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	ctx := context.Background()

	expectedSlug := "AVITO_TEST"
	expectedPercentage := uint8(10)

	querySelect := fmt.Sprintf(`
			SELECT deleted
			FROM %s
			WHERE slug = $1
	`, segmentsTable)

	queryInsert := fmt.Sprintf(`
			INSERT INTO %s (slug, auto_add_percentage)
			VALUES ($1, $2)
	`, segmentsTable)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelect)).WithArgs(expectedSlug).WillReturnError(pgx.ErrNoRows)
	mock.ExpectExec(regexp.QuoteMeta(queryInsert)).WithArgs(expectedSlug, expectedPercentage).
		WillReturnResult(pgxmock.NewResult("insert", 1))
	mock.ExpectCommit()

	storage := NewStoragePostgres()
	storage.db = mock

	err = storage.CreateSegment(ctx, expectedSlug, expectedPercentage)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestStorage_CreateSegmentExistedAndDeleted(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	ctx := context.Background()

	expectedSlug := "AVITO_TEST"
	expectedPercentage := uint8(10)

	querySelect := fmt.Sprintf(`
			SELECT deleted
			FROM %s
			WHERE slug = $1
	`, segmentsTable)

	queryUpdate := fmt.Sprintf(`
			UPDATE %s 
			SET deleted = false, 
				auto_add_percentage = $1
			WHERE slug = $2
	`, segmentsTable)

	expectedRow := pgxmock.NewRows([]string{"deleted"}).AddRow(true)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelect)).WithArgs(expectedSlug).
		WillReturnRows(expectedRow)
	mock.ExpectExec(regexp.QuoteMeta(queryUpdate)).WithArgs(expectedPercentage, expectedSlug).
		WillReturnResult(pgxmock.NewResult("update", 1))
	mock.ExpectCommit()

	storage := NewStoragePostgres()
	storage.db = mock

	err = storage.CreateSegment(ctx, expectedSlug, expectedPercentage)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestStorage_CreateSegmentExistedAndNotDeleted(t *testing.T) {
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

	querySelect := fmt.Sprintf(`
			SELECT deleted
			FROM %s
			WHERE slug = $1
	`, segmentsTable)

	expectedRow := pgxmock.NewRows([]string{"deleted"}).AddRow(false)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(querySelect)).WithArgs(expectedSlug).
		WillReturnRows(expectedRow)
	mock.ExpectCommit()

	storage := NewStoragePostgres()
	storage.db = mock

	err = storage.CreateSegment(ctx, expectedSlug, expectedPercentage)
	require.ErrorIs(t, err, expectedError)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}
