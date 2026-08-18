package db

import (
	"context"
	ty "danzmen/types"
	"database/sql"
	"fmt"
	"strings"
	"time"
)

// NOTE: tables
type DBMonthly_Task struct {
	Id         int
	Name       string
	Times_Done sql.NullFloat64
}

type DBMonthly_Record struct {
	Year_MonthId int
	MonthlyId    int
	Completed_At sql.NullString
}

type DBJoin_Monthly struct {
	*DBMonthly_Task
	*DBMonthly_Record
	*DBYear_Month

	ty.MonthlyTasksCfg
}

type DBLong_Tasks struct {
	Id           int
	Name         string
	Expires_in   sql.NullString
	Times_Done   sql.NullFloat64
	Completed_At sql.NullString

	ty.LongTermTasksCfg
}

type DBYear_Month struct {
	Id    int
	Year  int
	Month int
}

// NOTE: repository
func (s *SqliteDB) insertOrSelectYear_MonthID(date time.Time) (int, error) {
	ctx := context.Background()

	q := `
		insert into year_month (month_int, year) 
		values (?, ?) 
		on conflict(month_int, year) 
			do update set month_int = excluded.month_int
		returning id;
	`

	var id int
	err := s.db.QueryRowContext(ctx, q, int(date.Month()), date.Year()).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *SqliteDB) private_updateLogic(query string, mark_as_completed bool, args ...any) error {
	var c any = ""
	if mark_as_completed {
		c = ty.GetDate(ty.MM_DD_YYYY)
	}

	a := []any{c}
	for _, v := range args {
		a = append(a, v)
	}

	_, err := s.db.ExecContext(context.Background(), query, args...)
	return err
}

func (s *SqliteDB) UpdateCompletedMonthlyTask(taskid int, mark_as_completed bool) error {
	id, err := s.insertOrSelectYear_MonthID(time.Now())
	if err != nil {
		return err
	}

	q := `update monthly_record set completed_at = ? where monthly_id = ? AND year_month = ?;`
	return s.private_updateLogic(q, mark_as_completed, taskid, id)
}

func (s *SqliteDB) UpdateCompletedLongTask(taskid int, mark_as_completed bool) error {
	q := `update long_tasks set completed_at = ? where id = ?;`
	return s.private_updateLogic(q, mark_as_completed, taskid)
}

func (s *SqliteDB) CreateIfNotExistsMonthlyTasks(t []ty.MonthlyTasksCfg) ([]*DBJoin_Monthly, error) {
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
	q2 := `insert or ignore into monthly_record(monthly_id, year_month) values(?, ?);`

	ym_id, err := s.insertOrSelectYear_MonthID(time.Now())
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

		_, _ = tx.ExecContext(ctx, q2, t_id, ym_id)
		args = append(args, mcfg.Name) // so we save one more iteration
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	//NOTE: SELECT QUERY
	q3 := fmt.Sprintf(`
	select 
	t.id as t_id, t.name as t_name, t.times_done as t_times_done,
	ym.id as ym_id, ym.month_int as ym_month, ym.year as ym_year,
	d.year_month as d_year_month, d.monthly_id as d_monthlyid, d.completed_at as d_completed
	from monthly_tasks t
	left join monthly_record d on d.monthly_id = t.id and d.year_month = ?
	left join year_month ym on d.year_month = ym.id
	where t.name in (?%s) order by t.id asc;`, strings.Repeat(", ?", len(t)-1))

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
			&dt.Id, &dt.Name, &dt.Times_Done,
			&ym.Id, &ym.Month, &ym.Year,
			&dr.Year_MonthId, &dr.MonthlyId, &dr.Completed_At); err != nil {
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

func (s *SqliteDB) InsertOrSelectLongTermTasks(t []ty.LongTermTasksCfg) ([]*DBLong_Tasks, error) {
	if len(t) == 0 {
		return nil, fmt.Errorf("Not enough long term tasks")
	}

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	//then insert. ignore if they're dups.
	q1 := `insert or ignore into long_tasks(name, expires_in) values(?, ?);`

	cfgMap := make(map[string]ty.LongTermTasksCfg, len(t))
	n := []any{}

	for _, l := range t {
		if _, err = tx.ExecContext(ctx, q1, l.Name, l.MM_DD_YYYY_DATE); err != nil {
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
		`select id, name, expires_in, completed_at from long_tasks where name in (?%s)`,
		strings.Repeat(", ?", len(n)-1))

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

		if err := r.Scan(&t.Id, &t.Name, &t.Expires_in, &t.Completed_At); err != nil {
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
