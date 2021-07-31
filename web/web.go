package web

import (
	"net/http"

	"github.com/arken/arkstat/stats"
)

var global *stats.Stats

func Start(main *stats.Stats) {
	global = main
	// Setup handler functions for api endpoints
	http.HandleFunc("/usage", handleUsage)
	// Start http server and listen for incoming connections
	http.ListenAndServe(":8080", nil)
}
