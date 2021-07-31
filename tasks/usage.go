package tasks

import (
	"github.com/arken/arkstat/database"
	"github.com/arken/arkstat/stats"
)

func calculateUsage(db *database.DB, usage *stats.Usage) (err error) {
	usage.SpaceTotal, usage.SpaceUsed, err = db.GetPoolSize()
	if err != nil {
		return err
	}
	usage.NodesOnline, err = db.GetNodesOnline()
	return err
}
