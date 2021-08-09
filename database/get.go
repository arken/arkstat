package database

import (
	"database/sql"
	"time"
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
func (db *DB) GetAllOlderThan(input time.Duration) (result []Node, err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()
	// Ping database to check that it still exists.
	err = db.conn.Ping()
	if err != nil {
		return result, err
	}
	past := time.Now().Add(-1 * input)
	// Query for nodes that have not check in within the last 24 hours.
	rows, err := db.conn.Query("SELECT id, email FROM nodes WHERE last_seen < ?", past)
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

// GetPoolSize return the sum the nodes total space and used space.
func (db *DB) GetPoolSize() (total float64, used float64, err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()
	// Ping database to check that it still exists.
	err = db.conn.Ping()
	if err != nil {
		return total, used, err
	}
	// Get total pool size from sum of nodes reported values.
	row, err := db.conn.Query("SELECT SUM(total_space) FROM nodes")
	if err != nil {
		return total, used, err
	}
	defer row.Close()
	if !row.Next() {
		return total, used, sql.ErrNoRows
	}
	err = row.Scan(&total)
	if err != nil {
		return total, used, err
	}
	err = row.Close()
	if err != nil {
		return total, used, err
	}
	// Get used pool size from sum of nodes reported values.
	row, err = db.conn.Query("SELECT SUM(used_space) FROM nodes")
	if err != nil {
		return total, used, err
	}
	defer row.Close()
	if !row.Next() {
		return total, used, sql.ErrNoRows
	}
	err = row.Scan(&used)
	if err != nil {
		return total, used, err
	}
	return total, used, nil
}

// GetNodesOnline calculates the number of nodes reporting to the database.
func (db *DB) GetNodesOnline() (nodes int, err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()
	// Ping database to check that it still exists.
	err = db.conn.Ping()
	if err != nil {
		return -1, err
	}
	// Get total number of nodes reporting.
	row, err := db.conn.Query("SELECT COUNT(id) FROM nodes")
	if err != nil {
		return -1, err
	}
	defer row.Close()
	if !row.Next() {
		return nodes, sql.ErrNoRows
	}
	err = row.Scan(&nodes)
	if err != nil {
		return -1, err
	}
	return nodes, nil
}
