package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/custom_error"
)

func (s *Storage) CreateSegment(ctx context.Context, slug string, percentage uint8) error {
	tx, err := s.db.Begin(ctx)
	defer tx.Rollback(ctx)
	if err != nil {
		return err
	}

	querySelect := fmt.Sprintf(`
			SELECT deleted
			FROM %s
			WHERE slug = $1
	`, segmentsTable)

	var deleted bool

	err = tx.QueryRow(ctx, querySelect, slug).Scan(&deleted)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			queryInsert := fmt.Sprintf(`
					INSERT INTO %s (slug, auto_add_percentage)
					VALUES ($1, $2)
			`, segmentsTable)

			_, err = tx.Exec(ctx, queryInsert, slug, percentage)
			if err != nil {
				return err
			}

			err = tx.Commit(ctx)
			if err != nil {
				return err
			}

			return nil
		}
		return err
	}

	// if deleted is true, then segment existed before, and we change deleted to false
	if deleted {
		queryUpdate := fmt.Sprintf(`
				UPDATE %s 
				SET deleted = false, 
				    auto_add_percentage = $1
				WHERE slug = $2
		`, segmentsTable)

		_, err = tx.Exec(ctx, queryUpdate, percentage, slug)
		if err != nil {
			return err
		}

		err = tx.Commit(ctx)
		if err != nil {
			return err
		}

		return nil
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return custom_error.CustomError{
		Field:   "slug",
		Message: slug + " already exists",
	}
}
