package handler

import (
    "encoding/json"
    "net/http"
)

// hgetall input: /hgetall?key=myHash
func (h *Handler) HGetAll(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    key := r.URL.Query().Get("key")
    values, err := h.store.HGetAll(key)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]map[string]string{"value": values})
}
