package tasks

import (
	"math"

	"github.com/arken/arkstat/database"
	"github.com/arken/arkstat/stats"
)

func calculateUsage(db *database.DB, usage *stats.Usage) (err error) {
	total, used, err := db.GetPoolSize()
	if err != nil {
		return err
	}
	// Round Total and Used to 2 decimal places.
	usage.SpaceTotal = math.Round(total*100) / 100
	usage.SpaceUsed = math.Round(used*100) / 100
	usage.NodesOnline, err = db.GetNodesOnline()
	return err
}
