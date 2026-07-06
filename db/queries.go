package db

import (
	"context"
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
	q1 := `insert or ignore into daily_tasks(id, name) values(NULL, ?) returning id;`
	q2 := `insert into daily_record(daily_id) values(?);`

	for _, s := range names {
		var t_id int
		if err := tx.QueryRowContext(ctx, q1, s).Scan(&t_id); err != nil {
			continue
		}

		if _, err := tx.ExecContext(ctx, q2, t_id); err != nil {
			return nil, err
		}
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
	left join date_record d on d.task_id = t.id and d.date = date()
	where name in (?%s) order by t.id asc;`, strings.Repeat(", ?", len(names)-1))

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
