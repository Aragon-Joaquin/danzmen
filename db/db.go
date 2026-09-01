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
	CREATE TABLE IF NOT EXISTS monthly_tasks(
		id integer PRIMARY KEY,
		name text not null unique
	);

	CREATE TABLE IF NOT EXISTS long_tasks(
		id integer PRIMARY KEY,
		name text not null unique,
    expires_in text default(strftime('%m/%d/%Y', 'now', '+7 days')),
		times_done real,
		completed_at text null

		-- priority text not null check(priority IN ('low', 'med', 'high')) default('med'),
	);

	CREATE TABLE IF NOT EXISTS year_month(
		id integer primary key,
		month_int integer not null CHECK(month_int BETWEEN 1 AND 12),
		year integer not null CHECK(year > 2025),
		unique(month_int, year)
	);

	CREATE TABLE IF NOT EXISTS monthly_record(
		year_month integer not null,
		monthly_id integer not null,
		times_done real,

		completed_at text null,

		PRIMARY KEY (year_month, monthly_id),
		FOREIGN KEY(monthly_id) REFERENCES monthly_tasks(id)
	);

	CREATE TABLE IF NOT EXISTS monthly_progress(
		id integer PRIMARY KEY,
		date text not null default(date()) unique,
		tasks_completed int not null default 0
	);
	`)
	return err
}
