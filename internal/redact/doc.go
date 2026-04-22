// Package redact provides a Redactor type that masks sensitive configuration
// fields — such as passwords, tokens, and API keys — before they are written
// to reports, logs, or snapshots.
//
// Usage:
//
//	r := redact.New(nil)          // use default sensitive-key patterns
//	safe := r.Apply(configMap)    // returns a copy with sensitive values masked
//
// Custom patterns can be supplied to New to override the defaults:
//
//	r := redact.New([]string{"pin", "ssn"})
package redact
