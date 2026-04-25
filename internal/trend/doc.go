// Package trend provides drift-trend analysis for driftwatch services.
//
// It consumes history.Record slices produced by the history package and
// computes per-service metrics such as average diff count and drift direction
// (improving / stable / worsening). Results can be rendered as a table via
// Reporter.
package trend
