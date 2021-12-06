package web

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func handleUsage(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.Marshal(global.Usage)
	if err != nil {
		log.Printf("HTTP %d - %s", 500, err)
		http.Error(w, "Server Error", 500)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "max-age:290304000, public")
	w.Header().Set("Last-Modified", time.Now().Format(http.TimeFormat))
	w.Header().Set("Expires", time.Now().Add(5*time.Minute).Format(http.TimeFormat))
	w.Write(bytes)
}
