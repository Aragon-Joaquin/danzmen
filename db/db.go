package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

type SqliteDB struct {
	db *sql.DB
}

func Init() (*SqliteDB, error) {
	homePath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	dir := filepath.Join(homePath, ".local/share/danzmen")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", filepath.Join(dir, "danzmen.db"))
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
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS daily_tasks(
		id integer PRIMARY KEY,
		name text not null unique
	);

	CREATE TABLE IF NOT EXISTS long_tasks(
		id integer PRIMARY KEY,
		name text not null unique,
		expires_in text default(date('now','+7 days')),
		priority text not null check(priority IN ('low', 'med', 'high')) default('low'),
		completed_at text null
	);

	CREATE TABLE IF NOT EXISTS daily_record(
		date text not null default(date()),
		daily_id integer not null,
		completed integer not null check (completed IN (0, 1)) default (0),
		PRIMARY KEY (date, daily_id),
		FOREIGN KEY(daily_id) REFERENCES daily_tasks(id)
	);

	CREATE TABLE IF NOT EXISTS daily_progress(
		id integer PRIMARY KEY,
		date text not null default(date()) unique,
		tasks_completed int not null default 0
	);
	`)
	return err
}
