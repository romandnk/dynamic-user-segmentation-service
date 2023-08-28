package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/custom_error"
	"time"
)

func (s *Storage) CreateSegment(ctx context.Context, slug string, percentage uint8) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (slug, auto_add_percentage) 
		VALUES ($1, $2)
	`, segmentsTable)

	_, err := s.db.Exec(ctx, query, slug, percentage)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				return custom_error.CustomError{
					Field:   "slug",
					Message: slug + " already exists",
				}
			}
		}
		return fmt.Errorf("SegmentRepo.CreateSegment - s.db.Exec: %w", err)
	}

	return nil
}

func (s *Storage) DeleteSegment(ctx context.Context, slug string) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("SegmentRepo.DeleteSegment - s.db.Begin: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	queryDeleteFromUserSegments := fmt.Sprintf(`
		DELETE FROM %s
		WHERE segment_slug = $1
		RETURNING user_id
	`, userSegmentsTable)

	rows, err := tx.Query(ctx, queryDeleteFromUserSegments, slug)
	if err != nil {
		return fmt.Errorf("SegmentRepo.DeleteSegment - tx.Query: %w", err)
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var userID int

		err = rows.Scan(&userID)
		if err != nil {
			return fmt.Errorf("SegmentRepo.DeleteSegment - rows.Scan: %w", err)
		}

		userIDs = append(userIDs, userID)
	}

	queryDeleteFromSegments := fmt.Sprintf(`
		DELETE FROM %s
		WHERE slug = $1
	`, segmentsTable)

	ct, err := tx.Exec(ctx, queryDeleteFromSegments, slug)
	if err != nil {
		return fmt.Errorf("SegmentRepo.DeleteSegment - tx.Exec: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return custom_error.CustomError{
			Field:   "slug",
			Message: slug + " doesn't exist",
		}
	}

	queryAddOperations := fmt.Sprintf(`
		INSERT INTO %s (user_id, segment_slug, date, operation, auto_add) 
		VALUES ($1, $2, $3, $4, $5)
	`, operationsTable)

	now := time.Now().UTC()

	for _, userID := range userIDs {
		_, err = tx.Exec(ctx, queryAddOperations, userID, slug, now, "delete", false)
		if err != nil {
			return fmt.Errorf("SegmentRepo.DeleteSegment - tx.Exec: %w", err)
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("SegmentRepo.DeleteSegment - tx.Commit: %w", err)
	}

	return nil
}
