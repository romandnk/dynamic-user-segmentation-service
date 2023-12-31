package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pashagolub/pgxmock/v2"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/custom_error"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestStorage_UpdateUserSegments(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	ctx := context.Background()

	expectedSegmentsToAdd := []string{"AVITO_ADD"}
	expectedSegmentsToDelete := []string{"AVITO_DELETE"}
	expectedUserID := 1

	queryCheck := fmt.Sprintf(`
		SELECT true
		FROM %s
		WHERE user_id = $1 AND segment_slug = $2
	`, userSegmentsTable)

	queryInsertUserSegment := fmt.Sprintf(`
		INSERT INTO %s (user_id, segment_slug)
		VALUES ($1, $2)
	`, userSegmentsTable)

	queryInsertForAddOperation := fmt.Sprintf(`
		INSERT INTO %s (user_id, segment_slug, date, action, auto_add)
		VALUES ($1, $2, $3, $4, $5)
	`, operationsTable)

	queryDeleteUserSegment := fmt.Sprintf(`
		DELETE FROM %s
		WHERE user_id = $1 AND segment_slug = $2
	`, userSegmentsTable)

	queryInsertForDeleteOperation := fmt.Sprintf(`
		INSERT INTO %s (user_id, segment_slug, date, action, auto_add)
		VALUES ($1, $2, $3, $4, false)
	`, operationsTable)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).WithArgs(expectedUserID, expectedSegmentsToAdd[0]).
		WillReturnError(pgx.ErrNoRows)
	mock.ExpectExec(regexp.QuoteMeta(queryInsertUserSegment)).WithArgs(expectedUserID, expectedSegmentsToAdd[0]).
		WillReturnResult(pgxmock.NewResult("insert", 1))
	mock.ExpectExec(regexp.QuoteMeta(queryInsertForAddOperation)).
		WithArgs(expectedUserID, expectedSegmentsToAdd[0], pgxmock.AnyArg(), "add", false).
		WillReturnResult(pgxmock.NewResult("insert", 1))
	mock.ExpectExec(regexp.QuoteMeta(queryDeleteUserSegment)).WithArgs(expectedUserID, expectedSegmentsToDelete[0]).
		WillReturnResult(pgxmock.NewResult("delete", 1))
	mock.ExpectExec(regexp.QuoteMeta(queryInsertForDeleteOperation)).
		WithArgs(expectedUserID, expectedSegmentsToDelete[0], pgxmock.AnyArg(), "delete").
		WillReturnResult(pgxmock.NewResult("insert", 1))
	mock.ExpectCommit()

	storage := NewStoragePostgres()
	storage.db = mock

	err = storage.UpdateUserSegments(ctx, expectedSegmentsToAdd, expectedSegmentsToDelete, expectedUserID)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestStorage_UpdateUserSegmentsSegmentNotExistWhileAdding(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	ctx := context.Background()

	expectedSegmentsToAdd := []string{"AVITO_ADD"}
	expectedSegmentsToDelete := []string{"AVITO_DELETE"}
	expectedUserID := 1
	expectedError := custom_error.CustomError{
		Field:   "segments_to_add",
		Message: expectedSegmentsToAdd[0] + " doesn't exist",
	}

	queryCheck := fmt.Sprintf(`
		SELECT true
		FROM %s
		WHERE user_id = $1 AND segment_slug = $2
	`, userSegmentsTable)

	queryInsertUserSegment := fmt.Sprintf(`
		INSERT INTO %s (user_id, segment_slug)
		VALUES ($1, $2)
	`, userSegmentsTable)

	returnError := &pgconn.PgError{
		Code: "23503",
	}

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryCheck)).WithArgs(expectedUserID, expectedSegmentsToAdd[0]).
		WillReturnError(pgx.ErrNoRows)
	mock.ExpectExec(regexp.QuoteMeta(queryInsertUserSegment)).WithArgs(expectedUserID, expectedSegmentsToAdd[0]).
		WillReturnError(returnError)
	mock.ExpectRollback()

	storage := NewStoragePostgres()
	storage.db = mock

	err = storage.UpdateUserSegments(ctx, expectedSegmentsToAdd, expectedSegmentsToDelete, expectedUserID)
	require.ErrorIs(t, err, expectedError)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestStorage_UpdateUserSegmentsUserNotHaveSegmentWhileDeleting(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	ctx := context.Background()

	expectedSegmentsToAdd := []string{}
	expectedSegmentsToDelete := []string{"AVITO_DELETE"}
	expectedUserID := 1
	expectedError := custom_error.CustomError{
		Field:   "segments_to_delete",
		Message: fmt.Sprintf("User (%d) doesn't have segment %s", expectedUserID, expectedSegmentsToDelete[0]),
	}

	queryDeleteUserSegment := fmt.Sprintf(`
		DELETE FROM %s
		WHERE user_id = $1 AND segment_slug = $2
	`, userSegmentsTable)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryDeleteUserSegment)).WithArgs(expectedUserID, expectedSegmentsToDelete[0]).
		WillReturnResult(pgxmock.NewResult("insert", 0))
	mock.ExpectRollback()

	storage := NewStoragePostgres()
	storage.db = mock

	err = storage.UpdateUserSegments(ctx, expectedSegmentsToAdd, expectedSegmentsToDelete, expectedUserID)
	require.ErrorIs(t, err, expectedError)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestCheckUserSegment(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	ctx := context.Background()

	expectedSegment := "AVITO_TEST"
	expectedUserID := 1

	query := fmt.Sprintf(`
		SELECT true
		FROM %s
		WHERE user_id = $1 AND segment_slug = $2
	`, userSegmentsTable)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(expectedUserID, expectedSegment).
		WillReturnRows(pgxmock.NewRows([]string{"true"}).AddRow(true))

	err = checkUserSegment(ctx, mock, expectedSegment, expectedUserID)
	require.ErrorIs(t, err, ErrUserAlreadyHasSegment)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestStorage_GetActiveSegments(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	ctx := context.Background()

	expectedUserID := 1
	expectedUserSegments := []string{"TEST1", "TEST2"}

	query := fmt.Sprintf(`
		SELECT segment_slug
		FROM %s
		WHERE user_id = $1
	`, userSegmentsTable)

	mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(expectedUserID).
		WillReturnRows(pgxmock.NewRows([]string{"segment_slug"}).
			AddRow(expectedUserSegments[0]).
			AddRow(expectedUserSegments[1]))

	storage := NewStoragePostgres()
	storage.db = mock

	userSegments, err := storage.GetActiveSegments(ctx, expectedUserID)
	require.NoError(t, err)
	require.ElementsMatch(t, expectedUserSegments, userSegments)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}

func TestStorage_AutoAddUserSegments(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	ctx := context.Background()

	queryCount := fmt.Sprintf(`
		SELECT COUNT(DISTINCT user_id)
		FROM %s
	`, userSegmentsTable)

	querySelectSegments := fmt.Sprintf(`
		SELECT slug, auto_add_percentage
		FROM %s
		WHERE auto_add_percentage > 0
	`, segmentsTable)

	querySelectUsers := fmt.Sprintf(`
		SELECT COUNT(DISTINCT user_id)
		FROM %s
		WHERE segment_slug = $1
	`, userSegmentsTable)

	querySelectUsersWithoutCertainSegment := fmt.Sprintf(`
		SELECT DISTINCT user_id
		FROM %s
		WHERE user_id NOT IN (
    		SELECT DISTINCT user_id
    		FROM %s
    		WHERE segment_slug = $1
		)
		LIMIT $2
	`, userSegmentsTable, userSegmentsTable)

	queryInsertUserSegment := fmt.Sprintf(`
		INSERT INTO %s (user_id, segment_slug)
		VALUES ($1, $2)
	`, userSegmentsTable)

	queryInsertOperation := fmt.Sprintf(`
		INSERT INTO %s (user_id, segment_slug, date, action, auto_add)
		VALUES ($1, $2, $3, $4, $5)
	`, operationsTable)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryCount)).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(4))
	mock.ExpectQuery(regexp.QuoteMeta(querySelectSegments)).
		WillReturnRows(pgxmock.NewRows([]string{"slug", "auto_add_percentage"}).AddRow("TEST", 100))
	mock.ExpectQuery(regexp.QuoteMeta(querySelectUsers)).
		WithArgs("TEST").
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(1))
	mock.ExpectQuery(regexp.QuoteMeta(querySelectUsersWithoutCertainSegment)).WithArgs("TEST", 3).
		WillReturnRows(pgxmock.NewRows([]string{"user_id"}).
			AddRow(1).
			AddRow(2).
			AddRow(3))
	for i := 1; i <= 3; i++ {
		mock.ExpectExec(regexp.QuoteMeta(queryInsertUserSegment)).
			WithArgs(i, "TEST").
			WillReturnResult(pgxmock.NewResult("insert", 1))
		mock.ExpectExec(regexp.QuoteMeta(queryInsertOperation)).
			WithArgs(i, "TEST", pgxmock.AnyArg(), "add", true).
			WillReturnResult(pgxmock.NewResult("insert", 1))
	}
	mock.ExpectCommit()

	storage := NewStoragePostgres()
	storage.db = mock

	err = storage.AutoAddUserSegments(ctx)
	require.NoError(t, err)

	require.NoError(t, mock.ExpectationsWereMet(), "there was unexpected result")
}
