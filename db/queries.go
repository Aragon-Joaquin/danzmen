package db

import (
	"context"
	ty "danzmen/types"
	"database/sql"
	"fmt"
	"strings"
)

// NOTE: tables
type DBDaily_Task struct {
	Id   int
	Name string
}
type DBDaily_Record struct {
	Date      string
	DailyId   int
	Completed int
}

type DBJoin_Daily struct {
	*DBDaily_Task
	*DBDaily_Record
}

type DBLong_Tasks struct {
	Id           int
	Name         string
	Expires_in   sql.NullString
	Priority     ty.PRIORITY_TYPES
	Completed_at sql.NullString
}

// NOTE: repository
func (s *SqliteDB) InsertTask(name string) (*DBDaily_Task, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})

	defer tx.Rollback()

	if err != nil {
		return nil, err
	}

	q1 := `insert into daily_tasks(id, name) values (NULL, ?) returning id,name;`
	r := tx.QueryRowContext(ctx, q1, name)

	if err := r.Err(); err != nil {
		return nil, err
	}

	t := &DBDaily_Task{}
	if err := r.Scan(&t.Id, &t.Name); err != nil {
		return nil, err
	}

	q2 := `insert into daily_record(daily_id) values (?);`
	if _, err := tx.ExecContext(ctx, q2, t.Id); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return t, nil
}

func (s *SqliteDB) UpdateCompletedTask(id int, completed int) error {
	q := `
	update daily_record set completed = ? where daily_id = ? AND date = date();
	`

	_, err := s.db.ExecContext(context.Background(), q, completed, id)
	return err
}

func (s *SqliteDB) CreateIfNotExistsTasks(names []string) ([]*DBJoin_Daily, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	//NOTE: INSERT
	// create values or ignore errors
	// select the the values
	// insert them into daily_record
	q1 := `insert or ignore into daily_tasks(id, name) values(NULL, ?);`
	q1_5 := `select id from daily_tasks where name = (?)`
	q2 := `insert or ignore into daily_record(daily_id, date) values(?, date());`

	for _, s := range names {
		_, _ = tx.ExecContext(ctx, q1, s)

		var t_id int
		if err := tx.QueryRowContext(ctx, q1_5, s).Scan(&t_id); err != nil {
			return nil, err
		}

		_, _ = tx.ExecContext(ctx, q2, t_id)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	//NOTE: SELECT QUERY
	q3 := fmt.Sprintf(`
	select 
	t.id as t_id, t.name as t_name,
	coalesce(d.date, "") as d_date, d.daily_id as d_dailyid, d.completed as d_completed
	from daily_tasks t
	left join daily_record d on d.daily_id = t.id and d.date = date()
	where t.name in (?%s) order by t.id asc;`, strings.Repeat(", ?", len(names)-1))

	var args []any
	for _, n := range names {
		args = append(args, n)
	}

	r, err := s.db.QueryContext(ctx, q3, args...)
	if err != nil {
		return nil, err
	}

	res := []*DBJoin_Daily{}
	for r.Next() {
		//TODO: log this or return error?
		if r.Err() != nil {
			continue
		}

		t := &DBJoin_Daily{}
		dt := &DBDaily_Task{}
		dr := &DBDaily_Record{}

		if err := r.Scan(
			&dt.Id, &dt.Name,
			&dr.Date, &dr.DailyId, &dr.Completed); err != nil {
			return nil, err
		}

		t.DBDaily_Record = dr
		t.DBDaily_Task = dt
		res = append(res, t)
	}

	return res, nil
}

func (s *SqliteDB) InsertOrSelectLongTermTasks(tasks []ty.LongTermTasksCfg) ([]*DBLong_Tasks, error) {
	if len(tasks) == 0 {
		return nil, fmt.Errorf("Not enough long term tasks")
	}

	//validate if they're correctly parsed
	for _, t := range tasks {
		_, _, err := t.ValidateExpires_In()
		if err != nil {
			return nil, err
		}

		if err := t.ValidatePriority(); err != nil {
			return nil, err
		}
	}

	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	//then insert. ignore if they're dups.
	q1 := `insert or ignore into long_tasks(id, name, expires_in, priority) values(NULL, ?, ?, ?);`

	n := []any{}
	for _, t := range tasks {
		_, _ = tx.ExecContext(ctx, q1, t.Name, t.Ends, t.Priority)
		n = append(n, t.Name)
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return nil, err
	}

	//and select them
	q2 := fmt.Sprintf(
		`select id, name, expires_in, priority, completed_at from daily_tasks where name = (?%s)`,
		strings.Repeat(", ?", len(n)-1))

	r, err := s.db.QueryContext(ctx, q2, n...)
	if err != nil {
		return nil, err
	}

	DBTask := []*DBLong_Tasks{}
	for r.Next() {
		if r.Err() != nil {
			continue
		}
		t := DBLong_Tasks{}

		if err := r.Scan(&t.Id, &t.Name, &t.Expires_in, &t.Priority, &t.Completed_at); err != nil {
			return nil, err
		}

		DBTask = append(DBTask, &t)
	}

	return DBTask, nil
}
