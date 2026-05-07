package replay_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/driftwatch/internal/history"
	"github.com/driftwatch/internal/replay"
)

func writeRecord(t *testing.T, dir, service string, ts time.Time, diffs []string) {
	t.Helper()
	err := history.Record(dir, history.Entry{
		Service:   service,
		Timestamp: ts,
		Diffs:     diffs,
	})
	if err != nil {
		t.Fatalf("writeRecord: %v", err)
	}
}

func TestRun_NoEntries(t *testing.T) {
	dir := t.TempDir()
	r := replay.NewWithWriter(dir, io.Discard)

	from := time.Now().Add(-time.Hour)
	to := time.Now()
	res, err := r.Run("svc", from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 0 {
		t.Errorf("expected 0 entries, got %d", len(res.Entries))
	}
	if res.HasDrift() {
		t.Error("expected no drift")
	}
}

func TestRun_FiltersByWindow(t *testing.T) {
	dir := t.TempDir()
	now := time.Now().UTC().Truncate(time.Second)

	writeRecord(t, dir, "svc", now.Add(-2*time.Hour), []string{"a=1"})
	writeRecord(t, dir, "svc", now.Add(-30*time.Minute), []string{"b=2"})
	writeRecord(t, dir, "svc", now.Add(time.Hour), []string{"c=3"})

	var buf bytes.Buffer
	r := replay.NewWithWriter(dir, &buf)

	from := now.Add(-time.Hour)
	to := now
	res, err := r.Run("svc", from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(res.Entries))
	}
	if res.Entries[0].Diffs[0] != "b=2" {
		t.Errorf("unexpected diff: %s", res.Entries[0].Diffs[0])
	}
	if !res.HasDrift() {
		t.Error("expected drift")
	}
}

func TestReporter_Write_NoDrift(t *testing.T) {
	var buf bytes.Buffer
	rep := replay.NewReporter(&buf)
	res := replay.Result{
		Service: "svc",
		From:    time.Now().Add(-time.Hour),
		To:      time.Now(),
	}
	rep.Write(res)
	if !bytes.Contains(buf.Bytes(), []byte("No drift")) {
		t.Errorf("expected no-drift message, got: %s", buf.String())
	}
}

func TestReporter_Write_WithDrift(t *testing.T) {
	var buf bytes.Buffer
	rep := replay.NewReporter(&buf)
	res := replay.Result{
		Service: "svc",
		From:    time.Now().Add(-time.Hour),
		To:      time.Now(),
		Entries: []replay.Entry{
			{Service: "svc", Timestamp: time.Now(), Diffs: []string{"env=prod", "replicas=3"}},
		},
	}
	rep.Write(res)
	if !bytes.Contains(buf.Bytes(), []byte("env=prod")) {
		t.Errorf("expected diff in output, got: %s", buf.String())
	}
}

func TestNewReporter_NilFallback(t *testing.T) {
	// should not panic
	rep := replay.NewReporter(nil)
	res := replay.Result{Service: "svc", From: time.Now(), To: time.Now()}
	rep.Write(res)
}

// Ensure New uses os.Stdout without panicking.
func TestNew_DoesNotPanic(t *testing.T) {
	dir := t.TempDir()
	_ = replay.New(dir)
}

// io.Discard shim for older Go compat (Go 1.16+)
var _ = os.DevNull
var _ = filepath.Join
