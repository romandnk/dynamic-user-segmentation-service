package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/custom_error"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/models"
	"math"
	"time"
)

var ErrUserAlreadyHasSegment = errors.New("user alredy has segment")

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
		err = checkUserSegment(ctx, tx, segment, userID)
		if err != nil {
			if errors.Is(err, ErrUserAlreadyHasSegment) {
				continue
			}
			return err
		}
		err = addUserSegment(ctx, tx, segment, userID, false, now)
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

func checkUserSegment(ctx context.Context, tx pgx.Tx, segment string, userID int) error {
	var existFlag bool

	query := fmt.Sprintf(`
		SELECT true
		FROM %s
		WHERE user_id = $1 AND segment_slug = $2
	`, userSegmentsTable)

	err := tx.QueryRow(ctx, query, userID, segment).Scan(&existFlag)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return fmt.Errorf("UserRepo.checkUserSegment - tx.QueryRow.Scan: %w", err)
	}

	return ErrUserAlreadyHasSegment
}

func addUserSegment(ctx context.Context, tx pgx.Tx, segment string, userID int, autoAdd bool, now time.Time) error {
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
					Field:   "segments_to_add",
					Message: segment + " doesn't exist",
				}
			}
		}
		return fmt.Errorf("UserRepo.addUserSegment - tx.Exec: %w", err)
	}

	queryInsertOperation := fmt.Sprintf(`
		INSERT INTO %s (user_id, segment_slug, date, action, auto_add)
		VALUES ($1, $2, $3, $4, $5)
	`, operationsTable)

	_, err = tx.Exec(ctx, queryInsertOperation, userID, segment, now, "add", autoAdd)
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
			Field:   "segments_to_delete",
			Message: fmt.Sprintf("User (%d) doesn't have segment %s", userID, segment),
		}
	}

	queryInsertOperation := fmt.Sprintf(`
		INSERT INTO %s (user_id, segment_slug, date, action, auto_add)
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

func (s *Storage) AutoAddUserSegments(ctx context.Context) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("UserRepo.AutoAddUserSegments - s.db.Begin: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	now := time.Now().UTC()

	count, err := countUsers(ctx, tx)
	if err != nil {
		return err
	}

	querySelectSegments := fmt.Sprintf(`
		SELECT slug, auto_add_percentage
		FROM %s
		WHERE auto_add_percentage > 0
	`, segmentsTable)

	rows, err := tx.Query(ctx, querySelectSegments)
	if err != nil {
		return fmt.Errorf("UserRepo.AutoAddUserSegments - tx.Query: %w", err)
	}
	defer rows.Close()

	var segments []models.Segment
	for rows.Next() {
		var segment models.Segment

		err = rows.Scan(&segment.Slug, &segment.Percentage)
		if err != nil {
			return fmt.Errorf("UserRepo.AutoAddUserSegments - rows.Scan: %w", err)
		}

		segments = append(segments, segment)
	}

	for _, segment := range segments {
		numUsersToAdd := math.Floor(float64(count) * float64(segment.Percentage) / 100)
		err = addSegmentToUsers(ctx, tx, segment.Slug, int(numUsersToAdd), now)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("UserRepo.AutoAddUserSegments - tx.Commit: %w", err)
	}

	return nil
}

func countUsers(ctx context.Context, tx pgx.Tx) (int, error) {
	var count int

	query := fmt.Sprintf(`
		SELECT COUNT(DISTINCT user_id)
		FROM %s
	`, userSegmentsTable)

	err := tx.QueryRow(ctx, query).Scan(&count)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return 0, fmt.Errorf("UserRepo.countUsers - tx.QueryRow.Scan: %w", err)
		}
	}

	return count, nil
}

func addSegmentToUsers(ctx context.Context, tx pgx.Tx, segment string, amount int, now time.Time) error {
	// select num of users that already have such a segment, if that num >= percentage of segment then do nothing
	var num int

	querySelectUsers := fmt.Sprintf(`
		SELECT COUNT(DISTINCT user_id)
		FROM %s
		WHERE segment_slug = $1
	`, userSegmentsTable)

	err := tx.QueryRow(ctx, querySelectUsers, segment).Scan(&num)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("UserRepo.addSegmentToUsers - tx.QueryRow.Scan: %w", err)
		}
	}

	if num >= amount {
		return nil
	}

	// calculate remaining num of users that have to be with such a segment
	remainToAdd := amount - num

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

	rows, err := tx.Query(ctx, querySelectUsersWithoutCertainSegment, segment, remainToAdd)
	if err != nil {
		return fmt.Errorf("UserRepo.addSegmentToUsers - tx.Query: %w", err)
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var userID int

		err = rows.Scan(&userID)
		if err != nil {
			return fmt.Errorf("UserRepo.addSegmentToUsers - rows.Scan: %w", err)
		}

		userIDs = append(userIDs, userID)
	}

	for _, userID := range userIDs {
		err = addUserSegment(ctx, tx, segment, userID, true, now)
		if err != nil {
			return err
		}
	}

	return nil
}
