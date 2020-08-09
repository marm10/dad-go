package handlerException

import (
	"encoding/json"
	"net/http"
	"time"
)

type defaultError struct {
	StatusCode int       `json:"statusCode"`
	Message    string    `json:"mensagem"`
	TimeError  time.Time `json:"timestamp"`
	Path       string    `json:"path"`
}

func Handler(w http.ResponseWriter, r *http.Request, status int, message string) {
	de := defaultError{
		Message:    message,
		StatusCode: status,
		TimeError:  time.Now(),
		Path:       r.RequestURI,
	}
	w.WriteHeader(de.StatusCode)
	json.NewEncoder(w).Encode(&de)
}
