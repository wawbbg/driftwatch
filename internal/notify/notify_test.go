package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/driftwatch/internal/notify"
)

func serve(t *testing.T, status int, got *notify.Payload) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got != nil {
			if err := json.NewDecoder(r.Body).Decode(got); err != nil {
				t.Errorf("decode body: %v", err)
			}
		}
		w.WriteHeader(status)
	}))
}

func TestSend_Success(t *testing.T) {
	var received notify.Payload
	srv := serve(t, http.StatusOK, &received)
	defer srv.Close()

	n := notify.New(srv.URL)
	p := notify.Payload{
		Service:    "api-gateway",
		DriftCount: 2,
		Fields:     []string{"replicas", "image"},
		Timestamp:  time.Now().UTC(),
	}
	if err := n.Send(p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Service != "api-gateway" {
		t.Errorf("service = %q; want %q", received.Service, "api-gateway")
	}
	if received.DriftCount != 2 {
		t.Errorf("drift_count = %d; want 2", received.DriftCount)
	}
	if len(received.Fields) != 2 {
		t.Errorf("fields len = %d; want 2", len(received.Fields))
	}
}

func TestSend_NonOKStatus(t *testing.T) {
	srv := serve(t, http.StatusInternalServerError, nil)
	defer srv.Close()

	n := notify.New(srv.URL)
	err := n.Send(notify.Payload{Service: "svc"})
	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestSend_BadURL(t *testing.T) {
	n := notify.New("http://127.0.0.1:0/no-server")
	err := n.Send(notify.Payload{Service: "svc"})
	if err == nil {
		t.Fatal("expected error for unreachable URL, got nil")
	}
}

func TestSend_TimestampDefaulted(t *testing.T) {
	var received notify.Payload
	srv := serve(t, http.StatusNoContent, &received)
	defer srv.Close()

	n := notify.New(srv.URL)
	before := time.Now().UTC()
	if err := n.Send(notify.Payload{Service: "svc"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.Timestamp.Before(before) {
		t.Errorf("timestamp %v is before send time %v", received.Timestamp, before)
	}
}
