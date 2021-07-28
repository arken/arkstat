package database

import (
	"database/sql"
	"testing"
)

func TestGetRemoveEntry(t *testing.T) {
	// Initialize mock db
	db, err := openMock()
	if err != nil {
		t.Fatal(err)
	}
	// Create mock node
	in := Node{
		ID:         "im-not-a-real-node",
		Username:   "mrbaggins",
		Email:      "mrbaggins@example.com",
		TotalSpace: 4000,
		UsedSpace:  4,
	}
	// Add data to mock db
	err = db.Add(in)
	if err != nil {
		t.Fatal(err)
	}
	out, err := db.Get(in.ID)
	if err != nil {
		t.Fatal(err)
	}
	if out.ID != in.ID {
		t.Fatalf("In and Out IDs don't match!")
	}
	if out.Username != in.Username {
		t.Fatalf("In and Out Usernames don't match!")
	}
	if out.Email != in.Email {
		t.Fatalf("In and Out Emails don't match!")
	}
	if out.TotalSpace != in.TotalSpace {
		t.Fatalf("In and Out Total Space doesn't match!")
	}
	if out.UsedSpace != in.UsedSpace {
		t.Fatalf("In and Out Total Space doesn't match!")
	}
	// Remove Entry from DB
	err = db.Remove(in.ID)
	if err != nil {
		t.Fatal(err)
	}
	// Test that Get return nothing
	_, err = db.Get(in.ID)
	if err == nil || err != sql.ErrNoRows {
		t.Fatal("Test did not return a no rows error as expected")
	}

}
