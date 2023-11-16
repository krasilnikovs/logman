package handler

import (
	"encoding/json"
	"net/http"
)

// Index is a HandlerFunc which returns application info
func Index(w http.ResponseWriter, r *http.Request) {
	type ApplicationInfo struct {
		Application string `json:"application"`
	}

	data := ApplicationInfo{
		Application: "logman",
	}

	writeOkJson(w, data)
}

func writeOkJson(w http.ResponseWriter, data any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(data)
}
