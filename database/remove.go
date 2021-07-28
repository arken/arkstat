package database

// Remove deletes an node from the database.
func (db *DB) Remove(id string) (err error) {
	// Attempt to grab lock.
	db.lock.Lock()
	defer db.lock.Unlock()
	// Ping to check that database connection still exists.
	err = db.conn.Ping()
	if err != nil {
		return err
	}
	// Remove entry from database.
	stmt, err := db.conn.Prepare("DELETE FROM nodes WHERE id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	return err
}
