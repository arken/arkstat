package main

import (
	"log"

	"github.com/arken/arkstat/config"
	"github.com/arken/arkstat/database"
	"github.com/arken/arkstat/ipfs"
	"github.com/arken/arkstat/mail"
	"github.com/arken/arkstat/stats"
	"github.com/arken/arkstat/tasks"
	"github.com/arken/arkstat/web"
)

func main() {
	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	// Setup Database
	db, err := database.Open(config.Global.Database.Path)
	if err != nil {
		log.Fatal(err)
	}

	// Setup IPFS Node
	node, err := ipfs.CreateNode(config.Global.Ipfs.Path, ipfs.NodeConfArgs{
		Addr:           config.Global.Ipfs.Addr,
		PeerID:         config.Global.Ipfs.PeerID,
		PrivKey:        config.Global.Ipfs.PrivateKey,
		SwarmKey:       config.Manifest.ClusterKey,
		BootstrapPeers: config.Manifest.BootstrapPeers,
	})
	if err != nil {
		log.Fatal(err)
	}

	mailbox, err := mail.Init(config.Global.Mail.Domain,
		config.Global.Mail.PrivateKey,
		config.Global.Mail.Sender,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Setup Stats IPFS protocol handler
	node.SetHandler("/arkstat/0.0.1", ipfs.BuildStatsHandler(db, mailbox))

	// Initialize Stats Structure
	stats := stats.Stats{}

	// Startup background tasks
	go tasks.Start(db, &stats, mailbox)

	// Setup HTTP Server
	web.Start(config.Global.Web.Addr, &stats)
}
