package database

import (
	"database/sql"
)

// Add inserts a Node entry into the database if it doesn't exist already.
func (db *DB) Add(input Node) (err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()
	// Ping to check that database connection still exists.
	err = db.conn.Ping()
	if err != nil {
		return err
	}
	_, err = db.get(input.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			err = db.insert(input)
		} else {
			return err
		}
	} else {
		err = db.update(input)
	}
	return err
}

// Insert adds a Node entry to the database.
func (db *DB) insert(entry Node) (err error) {
	stmt, err := db.conn.Prepare(
		`INSERT INTO nodes(
			id,
			username,
			email,
			total_space,
			used_space,
			last_seen,
			first_seen
		) VALUES(?,?,?,?,?,datetime('now'),datetime('now'));`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		entry.ID,
		entry.Username,
		entry.Email,
		entry.TotalSpace,
		entry.UsedSpace,
	)
	return err
}

// update changes a file's status in the database.
func (db *DB) update(entry Node) (err error) {
	stmt, err := db.conn.Prepare(
		`UPDATE nodes SET
			username = ?,
			email = ?,
			total_space = ?,
			used_space = ?,
			last_seen = datetime('now')
			WHERE id = ?;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(
		entry.Username,
		entry.Email,
		entry.TotalSpace,
		entry.UsedSpace,
		entry.ID,
	)
	return err
}
