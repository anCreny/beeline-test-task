package http

import (
	"encoding/json"
	"main/internal/service"
	"net/http"
	"strings"
)

type WriteRequest struct {
	Message string `json:"message"`
}

// endpoint - [PUT] /queue/{topic}
func WriteToQueue(service service.QueueService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		var writeRequest WriteRequest

		path := strings.Split(r.URL.Path, "/")
		topicName := path[len(path)-1]

		if err := json.NewDecoder(r.Body).Decode(&writeRequest); err != nil || writeRequest.Message == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		msg, err := json.Marshal(writeRequest.Message)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err := service.WriteToQueue(topicName, msg); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
