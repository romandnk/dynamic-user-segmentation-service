package postgres

import (
	"context"
)

func (s *Storage) UpdateUserSegments(ctx context.Context, segmentsToAdd, segmentsToDelete []string, id int) error {
	//tx, err := s.db.Begin(ctx)
	//if err != nil {
	//	return err
	//}
	//defer tx.Rollback(ctx)
	//
	//for _, segment := range segmentsToAdd {
	//	// select percentage from segments where slug = segment and deleted = false
	//}
	//
	//pgx.Batch{}
	//tx.SendBatch()
	return nil
}
