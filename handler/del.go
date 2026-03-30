package handler

import (
	"encoding/json"
	"net/http"
)

// del input: /del?key=test1
func (h *Handler) Del(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodDelete {
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "key not be emty", http.StatusBadRequest)
			return
		}
		stats, err := h.store.Del(key)
		if !stats {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}
