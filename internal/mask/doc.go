// Package mask replaces sensitive field values in config maps with a fixed
// placeholder string before the data is written to any output (reports,
// exports, audit logs, etc.).
//
// Usage:
//
//	masker := mask.New()
//	safe := masker.Apply(rawConfig)
//
// Custom keys can be supplied via mask.NewWithKeys.
package mask
