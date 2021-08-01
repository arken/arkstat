package ipfs

import (
	"log"

	"github.com/arken/arkstat/database"
	"github.com/libp2p/go-libp2p-core/network"
	"github.com/polydawn/refmt/json"
)

type report struct {
	Email      string  `json:"email"`
	TotalSpace float64 `json:"total_space"`
	UsedSpace  float64 `json:"used_space"`
}

func BuildStatsHandler(db *database.DB) network.StreamHandler {
	return func(stream network.Stream) {
		defer stream.Close()
		// Create report and write read json input into struct
		r := report{}
		u := json.NewUnmarshaller(stream)
		err := u.Unmarshal(&r)
		if err != nil {
			return
		}
		if r.UsedSpace > r.TotalSpace || r.UsedSpace == 0 || r.TotalSpace == 0 {
			return
		}
		// Construct peer db entry
		peer := database.Node{
			ID:         stream.Conn().RemotePeer().Pretty(),
			Email:      r.Email,
			TotalSpace: r.TotalSpace,
			UsedSpace:  r.UsedSpace,
		}
		err = db.Add(peer)
		if err != nil {
			log.Printf("IPFS Handler <--> DB - %s", err)
		}
	}
}
