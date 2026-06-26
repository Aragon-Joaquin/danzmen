package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type DBComplete int

const (
	Completed_NO  DBComplete = 0
	Completed_YES DBComplete = 1
)

func (c DBComplete) IsCompleted() bool {
	return c == 1
}

func (c DBComplete) ToInt() int { return int(c) }

// NOTE: tables
type DBTask struct {
	Id   int
	Name string
}
type DBDate_Record struct {
	Date      string
	TaskId    int
	Completed DBComplete
}

type DBJoin_DateRecord_Tasks struct {
	*DBTask
	*DBDate_Record
}

// NOTE: repository
func (s *SqliteDB) InsertTask(name string) (*DBTask, error) {
	q := `insert into tasks(name) values (?) returning id,name;`
	r := s.db.QueryRowContext(context.Background(), q, name)

	if err := r.Err(); err != nil {
		return nil, err
	}

	t := &DBTask{}
	err := r.Scan(&t.Id, &t.Name)
	return t, err
}

func (s *SqliteDB) UpdateCompletedTask(id int, c DBComplete) error {
	q := `
	update date_record set completed = ?
	where task_id = ?;
	`

	_, err := s.db.ExecContext(context.Background(), q, c.ToInt(), id)
	return err
}

func (s *SqliteDB) CreateIfNotExistsTasks(names []string) ([]*DBJoin_DateRecord_Tasks, error) {
	ctx := context.Background()
	tx, err := s.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, err
	}

	//NOTE: INSERT
	for _, s := range names {
		q1 := `insert into tasks(name) values(?) on conflict do nothing;`
		if _, err := tx.ExecContext(ctx, q1, s); err != nil {
			//do something here... or maybe not
			continue
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	//NOTE: SELECT QUERY
	q2 := fmt.Sprintf(`
	select 
	t.id as t_id, t.name as t_name,
	d.date as d_date, d.task_id as d_taskid, d.completed as d_completed
	from tasks t
	left join date_record d on d.task_id = t.id
	where name in (?%s);`, strings.Repeat(", ?", len(names)-1))

	var args []any
	for range names {
		args = append(args, names)
	}

	r, err := s.db.QueryContext(ctx, q2, args...)

	res := make([]*DBJoin_DateRecord_Tasks, len(names))
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
