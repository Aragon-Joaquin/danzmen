package db

import (
	"context"
	ty "danzmen/types"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

func (s *SqliteDB) InsertOrSelectLongTermTasks(t []ty.LongTermTasksCfg, pageNumb int64) ([]*DBLong_Tasks, error) {
	if len(t) == 0 {
		return nil, fmt.Errorf("Not enough long term tasks")
	}

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	//then insert. update it if dup.
	//q1 := `insert or ignore into long_tasks(name, expires_in) values(?, ?);`
	q1 := `insert into long_tasks (name, expires_in, times_required) values (?, ?, ?)
					on conflict (name) do update set 
					times_required = excluded.times_required, 
					expires_in = excluded.expires_in;`

	cfgMap := make(map[string]ty.LongTermTasksCfg, len(t))
	n := []any{}

	for _, l := range t {
		if _, err = tx.ExecContext(ctx, q1, l.Name, l.MM_DD_YYYY_DATE, l.Times); err != nil {
			return nil, err
		}

		cfgMap[l.Name] = l
		n = append(n, l.Name)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	//and select them
	q2 := fmt.Sprintf(
		`select id, name, expires_in, completed_at, times_done, times_required from long_tasks where name in (?%s) limit ? offset ?;`,
		strings.Repeat(", ?", len(n)-1))

	//append limit + offset
	offset := s.calculate_offset(pageNumb, ty.AT_LEAST_NUMBER_OF_LONG_TASKS)
	n = append(n, ty.AT_LEAST_NUMBER_OF_LONG_TASKS, offset)

	r, err := s.db.QueryContext(ctx, q2, n...)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	DBTask := []*DBLong_Tasks{}
	for r.Next() {
		if err := r.Err(); err != nil {
			return nil, err
		}

		t := DBLong_Tasks{}

		if err := r.Scan(
			&t.Id, &t.Name, &t.Expires_in, &t.Completed_At,
			&t.Times_Done, &t.Times_Required); err != nil {
			return nil, err
		}

		a, ok := cfgMap[t.Name]
		if !ok {
			continue
		}

		t.LongTermTasksCfg = a
		DBTask = append(DBTask, &t)
	}

	return DBTask, nil
}

func (s *SqliteDB) InsertOrSelectMonthlyTasks(t []ty.MonthlyTasksCfg, pageNumb int64) ([]*DBJoin_Monthly, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	//NOTE: INSERT
	// create values or ignore errors
	// select the the values
	// insert them into monthly_record
	q1 := `insert into monthly_tasks(name) values(?) 
	       on conflict(name) do update set name=name
	       returning id;`
	//	q2 := `insert or ignore into monthly_record(monthly_id, year_month) values(?, ?) on conflict times_required = ? ;`

	q2 := `insert into monthly_record (monthly_id, year_month, times_required) values (?, ?, ?)
					on conflict (year_month, monthly_id) do update set times_required = excluded.times_required;`

	ym_id, err := s.insertOrSelectYear_MonthID(ctx, tx, time.Now())
	if err != nil {
		return nil, err
	}

	cfgMap := make(map[string]ty.MonthlyTasksCfg, len(t))

	var args []any // preparing args for the q3
	args = append(args, ym_id)

	for _, mcfg := range t {
		cfgMap[mcfg.Name] = mcfg

		var t_id int
		if err := tx.QueryRowContext(ctx, q1, mcfg.Name).Scan(&t_id); err != nil {
			return nil, err
		}

		_, _ = tx.ExecContext(ctx, q2, t_id, ym_id, mcfg.Times)
		args = append(args, mcfg.Name) // so we save one more iteration
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	//NOTE: SELECT QUERY
	q3 := fmt.Sprintf(`
	select 
	t.id as t_id, t.name as t_name,
	ym.id as ym_id, ym.month_int as ym_month, ym.year as ym_year,
	d.year_month as d_year_month, d.monthly_id as d_monthlyid, d.completed_at as d_completed,
	d.times_done as d_times_done, d.times_required as d_times_required
	from monthly_tasks t
	left join monthly_record d on d.monthly_id = t.id and d.year_month = ?
	left join year_month ym on d.year_month = ym.id
	where t.name in (?%s) 
	order by 
		case when d.completed_at is null then 0 else 1 end,
		t.id 
	asc limit ? offset ?;`, strings.Repeat(", ?", len(t)-1))

	//append limit + offset
	offset := s.calculate_offset(pageNumb, ty.AT_LEAST_NUMBER_OF_MONTHLY_TASKS)
	args = append(args, ty.AT_LEAST_NUMBER_OF_MONTHLY_TASKS, offset)

	r, err := s.db.QueryContext(ctx, q3, args...)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	res := []*DBJoin_Monthly{}
	for r.Next() {
		if err := r.Err(); err != nil {
			return nil, err
		}

		dt := &DBMonthly_Task{}
		dr := &DBMonthly_Record{}
		ym := &DBYear_Month{}

		if err := r.Scan(
			&dt.Id, &dt.Name,
			&ym.Id, &ym.Month, &ym.Year,
			&dr.Year_MonthId, &dr.MonthlyId, &dr.Completed_At, &dr.Times_Done, &dr.Times_Required); err != nil {
			return nil, err
		}

		cfg, ok := cfgMap[dt.Name]
		if !ok {
			continue
		}

		res = append(res, &DBJoin_Monthly{
			DBMonthly_Task:   dt,
			DBMonthly_Record: dr,
			DBYear_Month:     ym,
			MonthlyTasksCfg:  cfg,
		})
	}

	return res, nil
}

// WARN: private
func (s *SqliteDB) calculate_offset(pageNumb int64, numberTask ty.AT_LEAST_NUMBER) int64 {
	n := max(pageNumb, 1)
	return (n - 1) * int64(numberTask)
}

func (s *SqliteDB) insertOrSelectYear_MonthID(ctx context.Context, tx *sql.Tx, date time.Time) (int, error) {
	q := `
		insert into year_month (month_int, year) 
		values (?, ?) 
		on conflict(month_int, year) 
			do update set month_int = excluded.month_int
		returning id;
	`

	var id int
	err := tx.QueryRowContext(ctx, q, int(date.Month()), date.Year()).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
