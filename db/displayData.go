package db

import (
	"context"
	ty "danzmen/types"
	"database/sql"
	"time"
)

type DisplayData struct {
	CountLongTasks          int
	CountMonthlyTasks       int
	NearestLongTaskToExpire int //days, -1 if error/none
}

func (s *SqliteDB) GetDisplayData(ctx context.Context, now time.Time) (dd *DisplayData) {
	dd = &DisplayData{
		NearestLongTaskToExpire: -1,
	}
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})

	if err != nil {
		return
	}
	defer tx.Rollback()

	ymID, err := s.selectYearMonthIDTx(ctx, tx, now)
	if err != nil {
		return
	}

	qMonthly := `SELECT COUNT(*) FROM monthly_record WHERE year_month = ?`
	if err := tx.QueryRowContext(ctx, qMonthly, ymID).Scan(&dd.CountMonthlyTasks); err != nil {
		return
	}

	qLong := `SELECT COUNT(*) FROM long_tasks WHERE completed_at IS NULL`
	if err := tx.QueryRowContext(ctx, qLong).Scan(&dd.CountLongTasks); err != nil {
		return
	}

	qNearest := `SELECT expires_in FROM long_tasks WHERE completed_at IS NULL`
	rows, err := tx.QueryContext(ctx, qNearest)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var exp sql.NullString
		if err := rows.Scan(&exp); err != nil {
			return
		}

		if !exp.Valid || exp.String == "" {
			continue
		}

		t, err := time.Parse(string(ty.MM_DD_YYYY), exp.String)
		if err != nil {
			continue
		}

		if !t.After(now) {
			continue
		}

		days := int(t.Sub(now).Hours() / 24) // matches tui/list.go:344
		if dd.NearestLongTaskToExpire == -1 || days < dd.NearestLongTaskToExpire {
			dd.NearestLongTaskToExpire = days
		}
	}

	if err := rows.Err(); err != nil {
		return
	}

	if err := tx.Commit(); err != nil {
		return
	}

	return
}

// WARN: private
func (s *SqliteDB) selectYearMonthIDTx(ctx context.Context, tx *sql.Tx, now time.Time) (int, error) {
	q := `SELECT id FROM year_month WHERE month_int = ? AND year = ?`

	var id int
	err := tx.QueryRowContext(ctx, q, int(now.Month()), now.Year()).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
