// Package annotate provides persistent key/value annotations for drift results.
//
// Annotations are stored per-service as JSON files inside a configurable
// directory (e.g. .driftwatch/annotations/).  They are independent of
// snapshots and baselines and are intended to carry human-readable context
// such as ticket references, owner names, or suppression reasons.
//
// Usage:
//
//	// Add an annotation
//	err := annotate.Set(dir, "api-gateway", "ticket", "INFRA-42")
//
//	// Read all annotations for a service
//	anns, err := annotate.Get(dir, "api-gateway")
//
//	// Remove annotations by key
//	err = annotate.Delete(dir, "api-gateway", "ticket")
package annotate
