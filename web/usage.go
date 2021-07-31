package web

import (
	"encoding/json"
	"log"
	"net/http"
)

func handleUsage(w http.ResponseWriter, r *http.Request) {
	bytes, err := json.Marshal(global.Usage)
	if err != nil {
		log.Printf("HTTP %d - %s", 500, err)
		http.Error(w, "Server Error", 500)
	}
	w.Write(bytes)
}
