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
	if _, err := os.Stat(path); err != nil {
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

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}

// 		CREATE TABLE IF NOT EXISTS tasks (
// 			id STRING PRIMARY KEY,
// 			name TEXT NOT NULL,
// 			description TEXT DEFAULT "" NOT NULL,
// 			status TEXT DEFAULT "pending" NOT NULL,
// 			created_at DATETIME NOT NULL,
// 			updated_at DATETIME NOT NULL
// 		);

		// CREATE TABLE IF NOT EXISTS progress (
		// 	id STRING PRIMARY KEY,
		// 	task_id STRING NOT NULL,
		// 	status_init TEXT NOT NULL,
		// 	status_end TEXT,
		// 	created_at DATETIME NOT NULL,
		// 	updated_at DATETIME NOT NULL
		// 	finished_at DATETIME,
		// FOREIGN KEY(task_id) REFERENCES tasks(id)
		// );
