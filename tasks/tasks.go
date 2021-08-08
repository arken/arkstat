package tasks

import (
	"log"
	"time"

	"github.com/arken/arkstat/database"
	"github.com/arken/arkstat/mail"
	"github.com/arken/arkstat/stats"
)

func Start(db *database.DB, stats *stats.Stats, mailbox *mail.Mailbox) {
	ticker := time.NewTicker(5 * time.Minute)
	for {
		// Calculate the pool usage statistics.
		err := calculateUsage(db, &stats.Usage)
		if err != nil {
			log.Println(err)
		}

		// Clean up old nodes from the database.
		err = cleanUpOldNodes(db, mailbox)
		if err != nil {
			log.Println(err)
		}

		// Wait for ticker
		<-ticker.C
	}
}
