package ipfs

import (
	"encoding/json"
	"log"

	"github.com/arken/arkstat/database"
	"github.com/libp2p/go-libp2p-core/network"
)

type report struct {
	Email      string  `json:"email"`
	TotalSpace float64 `json:"total_space"`
	UsedSpace  float64 `json:"used_space"`
}

func BuildStatsHandler(db *database.DB) network.StreamHandler {
	return func(stream network.Stream) {
		defer stream.Close()
		// Create report and json input into struct
		r := report{}
		buf := make([]byte, 128)
		n, err := stream.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		err = json.Unmarshal(buf[:n], &r)
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
		log.Printf("Received Stats from Node: %s", stream.Conn().RemotePeer().Pretty())
	}
}
