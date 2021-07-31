package main

import (
	"log"

	"github.com/arken/arkstat/database"
	"github.com/arken/arkstat/stats"
	"github.com/arken/arkstat/tasks"
	"github.com/arken/arkstat/web"
)

func main() {
	// Setup Database
	db, err := database.Open("arkstat.db")
	if err != nil {
		log.Fatal(err)
	}
	// Initialize Stats Structure
	stats := stats.Stats{}
	// Startup background tasks
	go tasks.Start(db, &stats)
	// Setup HTTP Server
	web.Start(&stats)
}
