package web

import (
	"net/http"

	"github.com/arken/arkstat/stats"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var global *stats.Stats

func Start(addr string, main *stats.Stats) {
	global = main

	// Setup Chi Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://arken.github.io", "https://arken.io", "http://arken.io", "http://arken.github.io"},
		AllowedMethods: []string{"GET"},
	}))

	// Setup handler functions for api endpoints
	r.Get("/usage", handleUsage)

	// Start http server and listen for incoming connections
	http.ListenAndServe(addr, r)
}
