package database

import (
	"database/sql"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

func openMock() (result *DB, err error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}

	// Create the nodes table if is doesn't already exist.
	// This will also create the database if it doesn't exist.
	_, err = db.Exec(
		`CREATE TABLE IF NOT EXISTS nodes(
			id TEXT NOT NULL,
			username TEXT,
			email TEXT,
			total_space REAL,
			used_space REAL,
			last_seen DATETIME,
			first_seen DATETIME,
			
			PRIMARY KEY(id)
		);`,
	)

	result = &DB{
		conn: db,
		lock: sync.Mutex{},
	}

	return result, err

}
