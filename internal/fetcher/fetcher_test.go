package fetcher_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/driftwatch/internal/fetcher"
)

func serve(t *testing.T, status int, body interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
		if body != nil {
			_ = json.NewEncoder(w).Encode(body)
		}
	}))
}

func TestFetch_Success(t *testing.T) {
	srv := serve(t, 200, map[string]interface{}{"LOG_LEVEL": "info", "PORT": 8080})
	defer srv.Close()

	f := fetcher.New()
	state, err := f.Fetch("svc", srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state.Name != "svc" {
		t.Errorf("expected name svc, got %s", state.Name)
	}
	if state.Fields["LOG_LEVEL"] != "info" {
		t.Errorf("expected LOG_LEVEL=info, got %s", state.Fields["LOG_LEVEL"])
	}
	if state.Fields["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %s", state.Fields["PORT"])
	}
}

func TestFetch_NonOKStatus(t *testing.T) {
	srv := serve(t, 503, nil)
	defer srv.Close()

	f := fetcher.New()
	_, err := f.Fetch("svc", srv.URL)
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestFetch_InvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("not-json"))
	}))
	defer srv.Close()

	f := fetcher.New()
	_, err := f.Fetch("svc", srv.URL)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestFetch_BadURL(t *testing.T) {
	f := fetcher.New()
	_, err := f.Fetch("svc", "http://127.0.0.1:0/no-such-service")
	if err == nil {
		t.Fatal("expected error for unreachable URL")
	}
}
