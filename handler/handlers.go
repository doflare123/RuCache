package handler

import (
	storage "RuCache/Storage"
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
	mux.HandleFunc("/hget", h.HGet)
	mux.HandleFunc("/hgetall", h.HGetAll)
	mux.HandleFunc("/hset", h.HSet)
	mux.HandleFunc("/hdel", h.HDel)
}
