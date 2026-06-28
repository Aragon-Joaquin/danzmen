package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

// NOTE: tables
type DBTask struct {
	Id   int
	Name string
}
type DBDate_Record struct {
	Date      string
	TaskId    int
	Completed int
}

type DBJoin_DateRecord_Tasks struct {
	*DBTask
	*DBDate_Record
}

// NOTE: repository
func (s *SqliteDB) InsertTask(name string) (*DBTask, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})

	defer tx.Rollback()

	if err != nil {
		return nil, err
	}

	q1 := `insert into tasks(id, name) values (NULL, ?) returning id,name;`
	r := tx.QueryRowContext(ctx, q1, name)

	if err := r.Err(); err != nil {
		return nil, err
	}

	t := &DBTask{}
	if err := r.Scan(&t.Id, &t.Name); err != nil {
		return nil, err
	}

	q2 := `insert into date_record(task_id) values (?);`
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
	update date_record set completed = ? where task_id = ?;
	`

	_, err := s.db.ExecContext(context.Background(), q, completed, id)
	return err
}

func (s *SqliteDB) CreateIfNotExistsTasks(names []string) ([]*DBJoin_DateRecord_Tasks, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	//NOTE: INSERT
	q1 := `insert or ignore into tasks(id, name) values(NULL, ?) returning id;`
	q2 := `insert into date_record(task_id) values(?);`

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
	coalesce(d.date, "") as d_date, d.task_id as d_taskid, d.completed as d_completed
	from tasks t
	left join date_record d on d.task_id = t.id
	where name in (?%s) order by t.id asc;`, strings.Repeat(", ?", len(names)-1))

	var args []any
	for _, n := range names {
		args = append(args, n)
	}

	r, err := s.db.QueryContext(ctx, q3, args...)
	if err != nil {
		return nil, err
	}

	res := []*DBJoin_DateRecord_Tasks{}
	for r.Next() {
		t := &DBJoin_DateRecord_Tasks{}
		t.DBTask = &DBTask{}
		t.DBDate_Record = &DBDate_Record{}

		if err := r.Scan(
			&t.DBTask.Id, &t.DBTask.Name,
			&t.DBDate_Record.Date, &t.DBDate_Record.TaskId, &t.DBDate_Record.Completed); err != nil {
			return nil, err
		}
		res = append(res, t)
	}

	return res, nil
}
