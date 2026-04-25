package trend_test

import (
	"testing"
	"time"

	"github.com/example/driftwatch/internal/history"
	"github.com/example/driftwatch/internal/trend"
)

func rec(service string, diffs int, minutesAgo int) history.Record {
	return history.Record{
		Service:   service,
		DiffCount: diffs,
		Timestamp: time.Now().Add(-time.Duration(minutesAgo) * time.Minute),
	}
}

func TestAnalyse_Empty(t *testing.T) {
	tr := trend.Analyse(nil)
	if tr.Direction != trend.DirectionStable {
		t.Fatalf("expected stable, got %s", tr.Direction)
	}
}

func TestAnalyse_SingleRecord(t *testing.T) {
	tr := trend.Analyse([]history.Record{rec("svc-a", 3, 10)})
	if tr.Samples != 1 {
		t.Fatalf("expected 1 sample, got %d", tr.Samples)
	}
	if tr.AvgDiffs != 3.0 {
		t.Fatalf("expected avg 3.0, got %.1f", tr.AvgDiffs)
	}
	if tr.Direction != trend.DirectionStable {
		t.Fatalf("expected stable, got %s", tr.Direction)
	}
}

func TestAnalyse_Worsening(t *testing.T) {
	records := []history.Record{
		rec("svc-b", 1, 30),
		rec("svc-b", 5, 10),
	}
	tr := trend.Analyse(records)
	if tr.Direction != trend.DirectionWorsening {
		t.Fatalf("expected worsening, got %s", tr.Direction)
	}
}

func TestAnalyse_Improving(t *testing.T) {
	records := []history.Record{
		rec("svc-c", 5, 30),
		rec("svc-c", 1, 10),
	}
	tr := trend.Analyse(records)
	if tr.Direction != trend.DirectionImproving {
		t.Fatalf("expected improving, got %s", tr.Direction)
	}
}

func TestAnalyse_AvgDiffs(t *testing.T) {
	records := []history.Record{
		rec("svc-d", 2, 40),
		rec("svc-d", 4, 20),
		rec("svc-d", 6, 5),
	}
	tr := trend.Analyse(records)
	if tr.AvgDiffs != 4.0 {
		t.Fatalf("expected avg 4.0, got %.1f", tr.AvgDiffs)
	}
	if tr.Samples != 3 {
		t.Fatalf("expected 3 samples, got %d", tr.Samples)
	}
}

func TestServiceTrend_String(t *testing.T) {
	tr := trend.ServiceTrend{
		Service:   "svc-e",
		Samples:   2,
		AvgDiffs:  3.5,
		Direction: trend.DirectionStable,
		LastSeen:  time.Date(2024, 6, 1, 12, 0, 0, 0, time.UTC),
	}
	s := tr.String()
	if s == "" {
		t.Fatal("expected non-empty string")
	}
}
