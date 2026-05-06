// Package verify provides a high-level Verifier that loads a stored
// baseline for a named service and compares it against a live
// configuration map, returning a structured Result that describes
// any detected drift.
//
// Typical usage:
//
//	v := verify.New(".driftwatch/baselines")
//	res, err := v.Run("api-gateway", liveConfig)
//	if err != nil { ... }
//	if !res.OK { ... }
package verify
