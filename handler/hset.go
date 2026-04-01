package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

type HSetRequest struct {
	Key    string     `json:"key"`
	Fields [][]string `json:"fields"`
	TTL    *int64     `json:"ttl,omitempty"`
}

// hset input: json {"key": "myHash", "fields": [["field1","value1"], ["field2","value2"]], "ttl": 60}
func (h *Handler) HSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "content-type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	defer r.Body.Close()
	var req HSetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var ttl *time.Duration
	if req.TTL != nil {
		d := time.Duration(*req.TTL) * time.Second
		ttl = &d
	}

	stats, err := h.store.HSet(req.Key, req.Fields, ttl)
	if !stats {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
