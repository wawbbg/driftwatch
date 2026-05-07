package resolve_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/driftwatch/internal/config"
	"github.com/driftwatch/internal/fetcher"
	"github.com/driftwatch/internal/resolve"
)

func serve(t *testing.T, payload map[string]string, status int) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(payload)
	}))
}

func TestService_Success(t *testing.T) {
	srv := serve(t, map[string]string{"version": "1.2.3", "region": "us-east-1"}, 200)
	defer srv.Close()

	r := resolve.NewWithFetcher(fetcher.New())
	svc := config.Service{
		Name:     "api",
		URL:      srv.URL,
		Expected: map[string]string{"version": "1.2.3"},
	}

	res, err := r.Service(context.Background(), svc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ServiceName != "api" {
		t.Errorf("name = %q, want %q", res.ServiceName, "api")
	}
	if res.Live["version"] != "1.2.3" {
		t.Errorf("live version = %q, want %q", res.Live["version"], "1.2.3")
	}
	if res.Expected["version"] != "1.2.3" {
		t.Errorf("expected version = %q, want %q", res.Expected["version"], "1.2.3")
	}
}

func TestService_NoURL(t *testing.T) {
	r := resolve.New()
	_, err := r.Service(context.Background(), config.Service{Name: "empty"})
	if err == nil {
		t.Fatal("expected error for missing URL")
	}
}

func TestService_FetchError(t *testing.T) {
	r := resolve.New()
	svc := config.Service{Name: "bad", URL: "http://127.0.0.1:0/nope"}
	_, err := r.Service(context.Background(), svc)
	if err == nil {
		t.Fatal("expected fetch error")
	}
}

func TestAll_PartialErrors(t *testing.T) {
	srv := serve(t, map[string]string{"env": "prod"}, 200)
	defer srv.Close()

	r := resolve.NewWithFetcher(fetcher.New())
	cfg := &config.Config{
		Services: []config.Service{
			{Name: "ok", URL: srv.URL, Expected: map[string]string{"env": "prod"}},
			{Name: "broken", URL: "http://127.0.0.1:0/bad"},
		},
	}

	results, errs := r.All(context.Background(), cfg)
	if len(results) != 1 {
		t.Errorf("results = %d, want 1", len(results))
	}
	if len(errs) != 1 {
		t.Errorf("errs = %d, want 1", len(errs))
	}
}
