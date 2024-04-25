package db

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func Connect() (*Database, error) {
	path := "./chrono.db"
	newDB := false
	if _, err := os.Stat(path); err != nil {
		newDB = true
		file, err := os.Create(path)
		defer file.Close()
		if err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	d := &Database{db: db}
	if newDB {
		err = d.BuildSchema()
		if err != nil {
			return nil, err
		}
	}

	return d, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}

func (d *Database) BuildSchema() error {
	_, err := d.db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id STRING PRIMARY KEY,
			name TEXT NOT NULL,
			status TEXT DEFAULT "pending" NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		);
	`)
	if err != nil {
		return err
	}

	_, err = d.db.Exec(`
		CREATE TABLE IF NOT EXISTS progress (
			id STRING PRIMARY KEY,
			task_id STRING,
			status TEXT DEFAULT "in_progress" NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			finished_at DATETIME
		);
	`)
	if err != nil {
		return err
	}

	return nil
}
