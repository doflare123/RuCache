package handler

import (
	"encoding/json"
	"net/http"
)

// hget input: /hget?key=myHash&field=field1
func (h *Handler) HGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	field := r.URL.Query().Get("field")
	if key == "" || field == "" {
		http.Error(w, "key or field not be emty", http.StatusBadRequest)
		return
	}
	value, err := h.store.HGet(key, field)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"value": *value})
}
