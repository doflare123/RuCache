package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

type SetRequest struct {
	Key   string
	Value string
	TTL   *int64 `json:"ttl,omitempty"`
}

// set input: json {"key": "myKey", "value": "myValue", "ttl": 60}
func (h *Handler) Set(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "content-type must be application/json", http.StatusUnsupportedMediaType)
			return
		}
		defer r.Body.Close()
		var req SetRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		var ttl *time.Duration
		if req.TTL != nil {
			d := time.Duration(*req.TTL) * time.Second
			ttl = &d
		}
		stats, err := h.store.Set(req.Key, req.Value, ttl)
		if !stats {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
