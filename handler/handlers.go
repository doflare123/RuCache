package handler

import (
	storage "RuCache/Storage"
	"fmt"
	"net/http"
)

type Handler struct {
	store          *storage.Storage
	isShuttingDown func() bool
}

func NewHandler(store *storage.Storage, isShuttingDown func() bool) *Handler {
	return &Handler{store: store,
		isShuttingDown: isShuttingDown,
	}
}

func (h *Handler) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("/health", h.Health)
	mux.HandleFunc("/get", h.Get)
	mux.HandleFunc("/set", h.Set)
	mux.HandleFunc("/del", h.Del)
}

// handler Health
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if h.isShuttingDown() {
		http.Error(w, "Service is shutting down", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "OK")
}
