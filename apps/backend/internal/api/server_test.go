package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kimmaze027/death-project/apps/backend/internal/store"
)

func newTestServer() http.Handler {
	memoryStore := store.NewMemoryStore()
	return NewServer(memoryStore).Handler()
}

func TestHealthz(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	newTestServer().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestPostEvent_AcceptsValidPayload(t *testing.T) {
	t.Parallel()

	payload := map[string]any{
		"event_id":    "evt-a",
		"device_id":   "watch-api-1",
		"event_type":  "alive",
		"occurred_at": time.Date(2026, 3, 3, 12, 0, 0, 0, time.UTC).Format(time.RFC3339),
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/v1/events", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	newTestServer().ServeHTTP(rr, req)

	if rr.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", rr.Code)
	}
}

func TestPostEvent_RejectsMissingFields(t *testing.T) {
	t.Parallel()

	payload := map[string]any{
		"device_id": "watch-api-2",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/v1/events", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	newTestServer().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestPostSnooze_RejectsInvalidDuration(t *testing.T) {
	t.Parallel()

	payload := map[string]any{
		"device_id":      "watch-api-3",
		"duration_hours": 4,
		"start_at":       time.Now().Format(time.RFC3339),
		"end_at":         time.Now().Add(4 * time.Hour).Format(time.RFC3339),
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/v1/snoozes", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	newTestServer().ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestLatestStatus_AfterEvent(t *testing.T) {
	t.Parallel()

	handler := newTestServer()

	eventPayload := map[string]any{
		"event_id":    "evt-latest",
		"device_id":   "watch-api-4",
		"event_type":  "skip_sleeping",
		"occurred_at": time.Date(2026, 3, 3, 13, 0, 0, 0, time.UTC).Format(time.RFC3339),
	}
	eventBody, _ := json.Marshal(eventPayload)
	eventReq := httptest.NewRequest(http.MethodPost, "/v1/events", bytes.NewReader(eventBody))
	eventRR := httptest.NewRecorder()
	handler.ServeHTTP(eventRR, eventReq)
	if eventRR.Code != http.StatusAccepted {
		t.Fatalf("expected 202 from post event, got %d", eventRR.Code)
	}

	statusReq := httptest.NewRequest(http.MethodGet, "/v1/devices/watch-api-4/latest-status", nil)
	statusRR := httptest.NewRecorder()
	handler.ServeHTTP(statusRR, statusReq)

	if statusRR.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", statusRR.Code)
	}

	var parsed map[string]any
	if err := json.Unmarshal(statusRR.Body.Bytes(), &parsed); err != nil {
		t.Fatalf("failed to parse response: %v", err)
	}
	if parsed["latest_event_type"] != "skip_sleeping" {
		t.Fatalf("unexpected latest_event_type: %v", parsed["latest_event_type"])
	}
	if parsed["is_sleeping"] != true {
		t.Fatalf("expected is_sleeping=true, got %v", parsed["is_sleeping"])
	}
}

func TestLatestStatus_NotFound(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/v1/devices/unknown/latest-status", nil)
	rr := httptest.NewRecorder()

	newTestServer().ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodGet, "/v1/events", nil)
	rr := httptest.NewRecorder()

	newTestServer().ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}
