// Package suppress manages temporary suppressions of known drift fields.
//
// A suppression silences alerts for a specific service field for a defined
// duration. This is useful when a drift is intentional or already being
// addressed and repeated notifications would be noisy.
//
// Suppressions are stored as JSON files under a configurable directory,
// keyed by service name. Expired entries are ignored automatically.
package suppress
