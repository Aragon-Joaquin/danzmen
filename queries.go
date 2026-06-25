package main

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type DBComplete int

const (
	completed_no  DBComplete = 0
	completed_yes DBComplete = 1
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

// NOTE: repository
func (s *SqliteDB) InsertTask(name string) (*DBTask, error) {
	q := `insert into tasks(name) values (?) returning id,name;`
	r := s.QueryRowContext(context.Background(), q, name)

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
	where id = ?;
	`

	_, err := s.ExecContext(context.Background(), q, c.ToInt(), id)
	return err
}

func (s *SqliteDB) CreateIfNotExistsTasks(names []string) ([]DBTask, error) {
	ctx := context.Background()
	tx, err := s.BeginTx(ctx, &sql.TxOptions{})
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
	select * from tasks 
	where name in (?%s) 
	returning id, name;`, strings.Repeat(", ?", len(names)-1))

	var args []any
	for i, v := range names {
		args[i] = v
	}

	r, err := s.QueryContext(ctx, q2, args...)

	res := make([]DBTask, len(names))
	for r.Next() {
		t := DBTask{}
		if err := r.Scan(&t.Id, &t.Name); err != nil {
			return nil, err
		}
		res = append(res, t)
	}

	return res, nil
}
