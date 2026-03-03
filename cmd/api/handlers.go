package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/SuleymanyanArkadi/eventhub/internal/store"
	"github.com/SuleymanyanArkadi/eventhub/internal/task"
)

// makeHandlers создаёт маршруты, привязанные к переданному store.
func makeHandlers(s *store.MemoryStore) http.Handler {
	r := chi.NewRouter()

	// POST /v1/tasks
	r.Post("/v1/tasks", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Type    string `json:"type"`
			Payload string `json:"payload"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}
		if req.Type == "" {
			http.Error(w, "type required", http.StatusBadRequest)
			return
		}

		id := uuid.New().String()
		now := time.Now().UTC()
		t := &task.Task{
			ID:        id,
			Type:      req.Type,
			Payload:   req.Payload,
			Status:    task.StatusPending,
			CreatedAt: now,
			UpdatedAt: now,
		}
		if err := s.Create(t); err != nil {
			http.Error(w, "could not create", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(t)
	})

	// GET /v1/tasks/{id}
	r.Get("/v1/tasks/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		tk, err := s.Get(id)
		if err != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(tk)
	})

	// GET /v1/tasks?offset=0&limit=10
	r.Get("/v1/tasks", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		offset, _ := strconv.Atoi(q.Get("offset"))
		limit := 10
		if v := q.Get("limit"); v != "" {
			if li, err := strconv.Atoi(v); err == nil && li > 0 {
				limit = li
			}
		}
		list, _ := s.List(offset, limit)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(list)
	})

	return r
}
