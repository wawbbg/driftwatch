// Package resolve bridges the config and fetcher packages, turning a
// list of service definitions into paired (live, expected) config maps
// ready for drift detection.
//
// Usage:
//
//	r := resolve.New()
//	results, errs := r.All(ctx, cfg)
package resolve
