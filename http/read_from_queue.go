package http

import (
	"encoding/json"
	"errors"
	"main/internal/queue"
	"main/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type Response struct {
	Message string `json:"message"`
	Err     string `json:"err"`
}

// endpoint - [GET] /queue/{topic}
// queries:
// * timeout - int - used to overwrite default timeout setting
func ReadFromQueue(service service.QueueService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path, "/")
		topicName := path[len(path)-1]

		timeoutStr := r.URL.Query().Get("timeout")

		var timeout *int
		if timeoutStr != "" {
			t, err := strconv.Atoi(timeoutStr)
			if err != nil || t < 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			timeout = &t
		}

		var resp Response

		msg, err := service.ReadFromQueue(r.Context(), topicName, timeout)
		w.Header().Set("Content-Type", "application/json")
		if err != nil {
			if errors.Is(err, queue.ErrorNoMessage) {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}

			resp.Err = err.Error()
		} else {
			w.WriteHeader(http.StatusOK)
		}

		resp.Message = string(msg)
		json.NewEncoder(w).Encode(resp)
	}
}
