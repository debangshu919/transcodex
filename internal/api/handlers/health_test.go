package handlers_test

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/debangshu919/transcodex/internal/api/handlers"
)

func TestHealth(t *testing.T) {
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&logBuf, nil))

	handler := handlers.Health(logger)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %q", contentType)
	}

	expectedBody := `{"status": "ok"}`
	if rr.Body.String() != expectedBody {
		t.Errorf("expected body %q, got %q", expectedBody, rr.Body.String())
	}

	var resp struct {
		Status string `json:"status"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse JSON response: %v", err)
	}
	if resp.Status != "ok" {
		t.Errorf("expected status %q, got %q", "ok", resp.Status)
	}

	if !strings.Contains(logBuf.String(), "Health check") {
		t.Errorf("expected log output to contain 'Health check', got %q", logBuf.String())
	}
}

func TestHealth_Requests(t *testing.T) {
	tests := []struct {
		name   string
		method string
		target string
	}{
		{
			name:   "GET /health",
			method: http.MethodGet,
			target: "/health",
		},
		{
			name:   "GET /healthz",
			method: http.MethodGet,
			target: "/healthz",
		},
		{
			name:   "GET /health with query params",
			method: http.MethodGet,
			target: "/health?check=true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := slog.New(slog.DiscardHandler)
			handler := handlers.Health(logger)

			req := httptest.NewRequest(tt.method, tt.target, nil)
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != http.StatusOK {
				t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
			}
			if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("expected Content-Type application/json, got %q", ct)
			}
			if rr.Body.String() != `{"status": "ok"}` {
				t.Errorf("expected body %q, got %q", `{"status": "ok"}`, rr.Body.String())
			}
		})
	}
}

func TestHealth_Logging(t *testing.T) {
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	handler := handlers.Health(logger)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	var logEntry struct {
		Level string `json:"level"`
		Msg   string `json:"msg"`
	}
	if err := json.Unmarshal(logBuf.Bytes(), &logEntry); err != nil {
		t.Fatalf("failed to parse log JSON: %v", err)
	}

	if logEntry.Level != "INFO" {
		t.Errorf("expected log level INFO, got %q", logEntry.Level)
	}
	if logEntry.Msg != "Health check" {
		t.Errorf("expected log message 'Health check', got %q", logEntry.Msg)
	}
}

func TestHealth_LoggerLevelFiltered(t *testing.T) {
	var logBuf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&logBuf, &slog.HandlerOptions{
		Level: slog.LevelError,
	}))

	handler := handlers.Health(logger)
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if logBuf.Len() != 0 {
		t.Errorf("expected no log output when level is Error, got %q", logBuf.String())
	}
}

func TestHealth_DiscardHandler(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)
	handler := handlers.Health(logger)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
	if rr.Body.String() != `{"status": "ok"}` {
		t.Errorf("expected body %q, got %q", `{"status": "ok"}`, rr.Body.String())
	}
}
