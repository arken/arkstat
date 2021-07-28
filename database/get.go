package database

import (
	"database/sql"
)

// Get searches for and returns a the corresponding entry from the
// database if the entry exists.
func (db *DB) Get(id string) (result Node, err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()
	// Ping database to check that it still exists.
	err = db.conn.Ping()
	if err != nil {
		return result, err
	}
	return db.get(id)
}

func (db *DB) get(id string) (result Node, err error) {
	row, err := db.conn.Query("SELECT * FROM nodes WHERE id = ?", id)
	if err != nil {
		return result, err
	}
	defer row.Close()
	if !row.Next() {
		return result, sql.ErrNoRows
	}
	err = row.Scan(
		&result.ID,
		&result.Username,
		&result.Email,
		&result.TotalSpace,
		&result.UsedSpace,
		&result.LastSeen,
		&result.FirstSeen,
	)

	return result, err
}

// GetAllOld returns all of the old node entries in the database in a channel.
func (db *DB) GetAllOld() (result []Node, err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()
	// Ping database to check that it still exists.
	err = db.conn.Ping()
	if err != nil {
		return result, err
	}
	// Query for nodes that have not check in within the last 24 hours.
	rows, err := db.conn.Query("SELECT id, email FROM nodes WHERE last_seen < datetime('now', '-1 day')")
	if err != nil {
		return result, err
	}

	defer rows.Close()
	for rows.Next() {
		var node Node
		err = rows.Scan(
			&node.ID,
			&node.Email)
		if err != nil {
			return result, err
		}
		result = append(result, node)
	}
	return result, nil
}
