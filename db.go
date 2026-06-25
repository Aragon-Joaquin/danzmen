package main

import (
	"database/sql"
	"os"
	"path/filepath"
)

type SqliteDB struct {
	*sql.DB
}

func initDB() (*SqliteDB, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(homePath, ".local/share/danzmen")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", "./danzmen.db")
	if err != nil {
		return nil, err
	}

	//init transaction
	p := &SqliteDB{db}
	if err := p.createDatabase(); err != nil {
		return nil, err
	}

	return p, nil
}

func (s *SqliteDB) createDatabase() error {
	tx, err := s.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`
	CREATE TABLE IF NOT EXISTS tasks(
		id integer PRIMARY KEY autoincrement not null,
		name varchar(64) not null unique
	);

	CREATE TABLE IF NOT EXISTS date_record(
		date text not null default(date()),
		task_id integer not null,
		completed integer not null check (completed IN (0, 1)) default (0),
		PRIMARY KEY (date, task_id)
		FOREIGN KEY(task_id) REFERENCES tasks(id)
	);
	`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	return tx.Commit()
}
