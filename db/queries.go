package db

import (
	"context"
	ty "danzmen/types"
	"database/sql"
	"time"
)

func (s *SqliteDB) UpdateCompletedMonthlyTask(taskid int, mark_as_completed bool) error {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback()

	id, err := s.selectYearMonthIDTx(ctx, tx, time.Now())
	if err != nil {
		return err
	}

	q := `update monthly_record set completed_at = ? where monthly_id = ? AND year_month = ?;`
	return s.private_updateLogic(ctx, tx, q, mark_as_completed, taskid, id)
}

func (s *SqliteDB) UpdateCompletedLongTask(taskid int, mark_as_completed bool) error {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback()

	q := `update long_tasks set completed_at = ? where id = ?;`
	return s.private_updateLogic(ctx, tx, q, mark_as_completed, taskid)
}

func (s *SqliteDB) AddQuantityToMonthlyTask(taskid int, quantity float64) error {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback()

	id, err := s.selectYearMonthIDTx(ctx, tx, time.Now())
	if err != nil {
		return err
	}

	q := `update monthly_record set 
    times_done = coalesce(times_done, 0) + ?,
    completed_at = case
        when (coalesce(times_done, 0) + ?) >= times_required then ? else null
    end where monthly_id = ? and year_month = ?;`

	_, err = tx.ExecContext(context.Background(), q, quantity, quantity, ty.GetDate(ty.MM_DD_YYYY), taskid, id)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}
	return err
}

func (s *SqliteDB) AddQuantityToLongTask(taskid int, quantity float64) error {
	//TODO: this sums times_done twice
	q := `update long_tasks set 
    times_done = coalesce(times_done, 0) + ?,
    completed_at = case
        when (coalesce(times_done, 0) + ?) >= times_required then ? else null
    end where id = ?;`

	_, err := s.db.ExecContext(context.Background(), q, quantity, quantity, ty.GetDate(ty.MM_DD_YYYY), taskid)
	return err
}

// WARN: private
func (s *SqliteDB) private_updateLogic(ctx context.Context, tx *sql.Tx, query string, mark_as_completed bool, args ...any) error {
	var c any = nil
	if mark_as_completed {
		c = ty.GetDate(ty.MM_DD_YYYY)
	}

	a := []any{c}
	for _, v := range args {
		a = append(a, v)
	}

	if _, err := tx.ExecContext(ctx, query, a...); err != nil {
		return err
	}

	return tx.Commit()
}
