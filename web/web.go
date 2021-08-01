package web

import (
	"net/http"

	"github.com/arken/arkstat/stats"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var global *stats.Stats

func Start(addr string, main *stats.Stats) {
	global = main
	// Setup Chi Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	// Setup handler functions for api endpoints
	r.Get("/usage", handleUsage)
	// Start http server and listen for incoming connections
	http.ListenAndServe(addr, r)
}
