package database

import (
	"database/sql"
	"testing"
)

func TestGetNoEntry(t *testing.T) {
	// Initialize mock db
	db, err := openMock()
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Get("i-am-not-a-real-id")
	if err == nil || err != sql.ErrNoRows {
		t.Fatal("Test did not return a no rows error as expected")
	}
}

func TestGetEntry(t *testing.T) {
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
		t.Fatalf("expected node id with %s but got %s instead", in.ID, out.ID)
	}
	if out.Username != in.Username {
		t.Fatalf("expected node username with %s but got %s instead", in.Username, out.Username)
	}
	if out.Email != in.Email {
		t.Fatalf("expected node email with %s but got %s instead", in.Email, out.Email)
	}
	if out.TotalSpace != in.TotalSpace {
		t.Fatalf("expected node total space to be %f but got %f instead", in.TotalSpace, out.TotalSpace)
	}
	if out.UsedSpace != in.UsedSpace {
		t.Fatalf("expected node used space to be %f but got %f instead", in.UsedSpace, out.UsedSpace)
	}
}

func TestGetNodesOnline(t *testing.T) {
	// Initialize mock db
	db, err := openMock()
	if err != nil {
		t.Fatal(err)
	}
	// Create mock node
	one := Node{
		ID:         "im-not-a-real-node",
		Username:   "mrbaggins",
		Email:      "mrbaggins@example.com",
		TotalSpace: 4000,
		UsedSpace:  4,
	}
	// Create another mock node
	two := Node{
		ID:         "im-also-not-a-real-node",
		Username:   "frodo",
		Email:      "frodo@example.com",
		TotalSpace: 8000,
		UsedSpace:  8,
	}
	// Add node one to mock db
	err = db.Add(one)
	if err != nil {
		t.Fatal(err)
	}
	// Add node two to mock db
	err = db.Add(two)
	if err != nil {
		t.Fatal(err)
	}
	// Check nodes online
	online, err := db.GetNodesOnline()
	if err != nil {
		t.Fatal(err)
	}
	if online != 2 {
		t.Fatalf("expected %d nodes online but got %d instead", 2, online)
	}
}

func TestGetPoolSize(t *testing.T) {
	// Initialize mock db
	db, err := openMock()
	if err != nil {
		t.Fatal(err)
	}
	// Create mock node
	one := Node{
		ID:         "im-not-a-real-node",
		Username:   "mrbaggins",
		Email:      "mrbaggins@example.com",
		TotalSpace: 4000,
		UsedSpace:  4,
	}
	// Create another mock node
	two := Node{
		ID:         "im-also-not-a-real-node",
		Username:   "frodo",
		Email:      "frodo@example.com",
		TotalSpace: 8000,
		UsedSpace:  8,
	}
	// Add node one to mock db
	err = db.Add(one)
	if err != nil {
		t.Fatal(err)
	}
	// Add node two to mock db
	err = db.Add(two)
	if err != nil {
		t.Fatal(err)
	}
	// Check pool size values
	total, used, err := db.GetPoolSize()
	if err != nil {
		t.Fatal(err)
	}
	if total != one.TotalSpace+two.TotalSpace {
		t.Fatalf("expected a total pool size of %f but got %f instead", one.TotalSpace+two.TotalSpace, total)
	}
	if used != one.UsedSpace+two.UsedSpace {
		t.Fatalf("expected a used pool size of %f but got %f instead", one.UsedSpace+two.UsedSpace, used)
	}
}
