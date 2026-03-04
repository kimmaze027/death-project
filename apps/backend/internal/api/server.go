package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kimmaze027/death-project/apps/backend/internal/domain"
	"github.com/kimmaze027/death-project/apps/backend/internal/store"
)

type Server struct {
	store *store.MemoryStore
}

func NewServer(store *store.MemoryStore) *Server {
	return &Server{store: store}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", s.healthz)
	mux.HandleFunc("/v1/events", s.postEvent)
	mux.HandleFunc("/v1/snoozes", s.postSnooze)
	mux.HandleFunc("/v1/devices/", s.getLatestStatus)
	return loggingMiddleware(mux)
}

func (s *Server) healthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (s *Server) postEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}

	var event domain.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if event.DeviceID == "" || event.EventType == "" || event.OccurredAt.IsZero() {
		writeError(w, http.StatusBadRequest, errors.New("missing required event fields"))
		return
	}

	s.store.AddEvent(event)
	writeJSON(w, http.StatusAccepted, map[string]any{
		"accepted":    true,
		"server_time": time.Now().Format(time.RFC3339),
	})
}

func (s *Server) postSnooze(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}

	var snooze domain.Snooze
	if err := json.NewDecoder(r.Body).Decode(&snooze); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	if snooze.DeviceID == "" || snooze.DurationHours < 1 || snooze.DurationHours > 3 {
		writeError(w, http.StatusBadRequest, errors.New("invalid snooze payload"))
		return
	}

	s.store.SetSnooze(snooze)
	writeJSON(w, http.StatusAccepted, map[string]bool{"accepted": true})
}

func (s *Server) getLatestStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		return
	}

	if !strings.HasSuffix(r.URL.Path, "/latest-status") {
		writeError(w, http.StatusNotFound, errors.New("not found"))
		return
	}

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(parts) != 4 || parts[0] != "v1" || parts[1] != "devices" || parts[3] != "latest-status" {
		writeError(w, http.StatusNotFound, errors.New("not found"))
		return
	}
	deviceID := parts[2]

	status, err := s.store.LatestStatus(deviceID)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}
	writeJSON(w, http.StatusOK, status)
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]string{"error": err.Error()})
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
