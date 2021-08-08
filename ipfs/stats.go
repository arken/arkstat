package ipfs

import (
	"encoding/json"
	"log"

	"github.com/arken/arkstat/database"
	"github.com/arken/arkstat/mail"
	"github.com/libp2p/go-libp2p-core/network"
)

type report struct {
	Email      string  `json:"email"`
	TotalSpace float64 `json:"total_space"`
	UsedSpace  float64 `json:"used_space"`
}

func BuildStatsHandler(db *database.DB, mailbox *mail.Mailbox) network.StreamHandler {
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

		// Parse input json.
		err = json.Unmarshal(buf[:n], &r)
		if err != nil {
			return
		}

		// Check for mailformed stats
		if r.UsedSpace > r.TotalSpace || r.UsedSpace == 0 || r.TotalSpace == 0 {
			return
		}

		// Construct peer db entry
		peer := database.Node{
			ID:         stream.Conn().RemotePeer().Pretty(),
			Email:      r.Email,
			TotalSpace: r.TotalSpace / 1000,
			UsedSpace:  r.UsedSpace / 1000,
		}

		// Add info to database
		err = db.Add(peer)
		if err != nil {
			log.Printf("IPFS Handler <--> DB - %s", err)
		}
		log.Printf("Received Stats from Node: %s", stream.Conn().RemotePeer().Pretty())

		// Send confirmation email
		if peer.Email != "" && mailbox != nil {
			err = mailbox.Send("mail/templates/welcome.yml", peer.Email)
			log.Println(err)
		}
	}
}
