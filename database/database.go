package database

import (
	"database/sql"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3" // Import sqlite driver for database interaction.
)

type DB struct {
	conn *sql.DB
	lock sync.Mutex
}

// Node defines the columns of data stored in the database.
type Node struct {
	ID         string
	Username   string
	Email      string
	TotalSpace float64
	UsedSpace  float64
	LastSeen   time.Time
	FirstSeen  time.Time
}

// Open connects and returns a pointer to a database object.
func Open(path string) (result *DB, err error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return result, err
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
		); PRAGMA busy_timeout = 5000;`,
	)

	result = &DB{
		conn: db,
		lock: sync.Mutex{},
	}
	return result, err
}
