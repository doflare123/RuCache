package handler

import (
	"fmt"
	"net/http"
)

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if h.isShuttingDown() {
		http.Error(w, "Service is shutting down", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}
