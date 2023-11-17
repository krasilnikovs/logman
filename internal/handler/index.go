package handler

import (
	"encoding/json"
	"net/http"

	"github.com/krasilnikovm/logman/internal/service"
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

func writeValidationJson(w http.ResponseWriter, err service.ErrValidation) {
	type validationResponse struct {
		Errors []string `json:"errors"`
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)

	json.NewEncoder(w).Encode(validationResponse{Errors: err.Errors})
}

func writeWithEmptyBody(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}
