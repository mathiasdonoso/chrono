package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	db *sql.DB
}

func Connect() (*Database, error) {
	log.Println("Creating connection to db")
	path := "./chrono.db"
	if _, err := os.Stat(path); err != nil {
		file, err := os.Create(path)
		defer file.Close()
		if err != nil {
			return nil, err
		}
	}

	log.Println("DB file exists, trying to open it")

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		log.Println("connect to db error:", err)
		return nil, err
	}

	log.Println("DB connection created")

	return &Database{db: db}, nil
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}

// func CreateDB() error {
// 	tasks := `
// 		CREATE TABLE IF NOT EXISTS tasks (
// 			id STRING PRIMARY KEY,
// 			name TEXT NOT NULL,
// 			description TEXT DEFAULT "" NOT NULL,
// 			status TEXT DEFAULT "pending" NOT NULL,
// 			created_at DATETIME NOT NULL,
// 			updated_at DATETIME NOT NULL
// 		);
// 	`
// 	works := `
// 		CREATE TABLE IF NOT EXISTS works (
// 			id STRING PRIMARY KEY,
// 			task_id STRING,
// 			status_init TEXT NOT NULL,
// 			status_end TEXT,
// 			created_at DATETIME NOT NULL,
// 			updated_at DATETIME NOT NULL
// 			finished_at DATETIME
// 		);
// 	`
//
// 	return createTables(tasks, works)
// }
//
// func createTables(sql ...string) error {
// 	for _, s := range sql {
// 		statement, err := db.Prepare(s)
// 		if err != nil {
// 			return err
// 		}
// 		_, err = statement.Exec()
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

