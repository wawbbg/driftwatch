//go:build integration

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

// TestAll_AllSucceed verifies that All resolves multiple healthy services
// without errors and returns one result per service.
func TestAll_AllSucceed(t *testing.T) {
	make := func(payload map[string]string) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_ = json.NewEncoder(w).Encode(payload)
		}))
	}

	s1 := make(map[string]string{"version": "2.0.0"})
	s2 := make(map[string]string{"version": "3.1.0"})
	defer s1.Close()
	defer s2.Close()

	cfg := &config.Config{
		Services: []config.Service{
			{Name: "svc-a", URL: s1.URL, Expected: map[string]string{"version": "2.0.0"}},
			{Name: "svc-b", URL: s2.URL, Expected: map[string]string{"version": "3.1.0"}},
		},
	}

	r := resolve.NewWithFetcher(fetcher.New())
	results, errs := r.All(context.Background(), cfg)

	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(results) != 2 {
		t.Fatalf("results = %d, want 2", len(results))
	}
	for _, res := range results {
		if res.Live["version"] == "" {
			t.Errorf("service %q: live version empty", res.ServiceName)
		}
	}
}
