package tasks

import (
	"log"
	"time"

	"github.com/arken/arkstat/database"
	"github.com/arken/arkstat/stats"
)

func Start(db *database.DB, stats *stats.Stats) {
	for range time.Tick(time.Minute * 5) {
		// Calculate the pool usage statistics.
		err := calculateUsage(db, &stats.Usage)
		if err != nil {
			log.Println(err)
		}
	}
}
