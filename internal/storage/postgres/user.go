package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/custom_error"
	"time"
)

func (s *Storage) UpdateUserSegments(ctx context.Context, segmentsToAdd, segmentsToDelete []string, userID int) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("UserRepo.UpdateUserSegments - s.db.Begin: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	now := time.Now().UTC()

	for _, segment := range segmentsToAdd {
		err = addUserSegment(ctx, tx, segment, userID, now)
		if err != nil {
			return err
		}
	}

	for _, segment := range segmentsToDelete {
		err = deleteUserSegment(ctx, tx, segment, userID, now)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("UserRepo.UpdateUserSegments - tx.Commit: %w", err)
	}

	return nil
}

func addUserSegment(ctx context.Context, tx pgx.Tx, segment string, userID int, now time.Time) error {
	queryInsertUserSegment := fmt.Sprintf(`
		INSERT INTO %s (user_id, segment_slug)
		VALUES ($1, $2)
	`, userSegmentsTable)

	_, err := tx.Exec(ctx, queryInsertUserSegment, userID, segment)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23503" {
				return custom_error.CustomError{
					Field:   "slug",
					Message: segment + " doesn't exist",
				}
			}
		}
		return fmt.Errorf("UserRepo.addUserSegment - tx.Exec: %w", err)
	}

	queryInsertOperation := fmt.Sprintf(`
		INSERT INTO %s (user_id, segment_slug, date, operation, auto_add)
		VALUES ($1, $2, $3, $4, false)
	`, operationsTable)

	_, err = tx.Exec(ctx, queryInsertOperation, userID, segment, now, "add")
	if err != nil {
		return fmt.Errorf("UserRepo.addUserSegment - tx.Exec: %w", err)
	}

	return nil
}

func deleteUserSegment(ctx context.Context, tx pgx.Tx, segment string, userID int, now time.Time) error {
	queryDeleteUserSegment := fmt.Sprintf(`
		DELETE FROM %s
		WHERE user_id = $1 AND segment_slug = $2
	`, userSegmentsTable)

	ct, err := tx.Exec(ctx, queryDeleteUserSegment, userID, segment)
	if err != nil {
		return fmt.Errorf("UserRepo.deleteUserSegment - tx.Exec: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return custom_error.CustomError{
			Field:   "slug",
			Message: fmt.Sprintf("User (%d) doesn't have segment %s", userID, segment),
		}
	}

	queryInsertOperation := fmt.Sprintf(`
		INSERT INTO %s (user_id, segment_slug, date, operation, auto_add)
		VALUES ($1, $2, $3, $4, false)
	`, operationsTable)

	_, err = tx.Exec(ctx, queryInsertOperation, userID, segment, now, "delete")
	if err != nil {
		return fmt.Errorf("UserRepo.deleteUserSegment - tx.Exec: %w", err)
	}

	return nil
}

func (s *Storage) GetActiveSegments(ctx context.Context, userID int) ([]string, error) {
	query := fmt.Sprintf(`
		SELECT segment_slug
		FROM %s
		WHERE user_id = $1
	`, userSegmentsTable)

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("UserRepo.GetActiveSegments - s.db.Query: %w", err)
	}
	defer rows.Close()

	var segments []string
	for rows.Next() {
		var segment string

		err = rows.Scan(&segment)
		if err != nil {
			return nil, fmt.Errorf("UserRepo.GetActiveSegments - rows.Scan: %w", err)
		}

		segments = append(segments, segment)
	}

	return segments, nil
}
