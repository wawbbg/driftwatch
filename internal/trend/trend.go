// Package trend analyses drift history to surface recurring and worsening
// drift patterns across services over time.
package trend

import (
	"fmt"
	"sort"
	"time"

	"github.com/example/driftwatch/internal/history"
)

// Direction describes whether drift for a service is improving or worsening.
type Direction string

const (
	DirectionStable   Direction = "stable"
	DirectionWorsening Direction = "worsening"
	DirectionImproving Direction = "improving"
)

// ServiceTrend summarises drift behaviour for a single service.
type ServiceTrend struct {
	Service    string
	Samples    int
	AvgDiffs   float64
	Direction  Direction
	LastSeen   time.Time
}

func (t ServiceTrend) String() string {
	return fmt.Sprintf("%s: samples=%d avg_diffs=%.1f direction=%s last=%s",
		t.Service, t.Samples, t.AvgDiffs, t.Direction, t.LastSeen.Format(time.RFC3339))
}

// Analyse reads history records for the given service and returns a
// ServiceTrend describing its drift behaviour.
func Analyse(records []history.Record) ServiceTrend {
	if len(records) == 0 {
		return ServiceTrend{Direction: DirectionStable}
	}

	sort.Slice(records, func(i, j int) bool {
		return records[i].Timestamp.Before(records[j].Timestamp)
	})

	var total int
	for _, r := range records {
		total += r.DiffCount
	}

	avg := float64(total) / float64(len(records))
	last := records[len(records)-1]

	dir := DirectionStable
	if len(records) >= 2 {
		recent := records[len(records)-1].DiffCount
		prior := records[len(records)-2].DiffCount
		switch {
		case recent > prior:
			dir = DirectionWorsening
		case recent < prior:
			dir = DirectionImproving
		}
	}

	return ServiceTrend{
		Service:   last.Service,
		Samples:   len(records),
		AvgDiffs:  avg,
		Direction: dir,
		LastSeen:  last.Timestamp,
	}
}
