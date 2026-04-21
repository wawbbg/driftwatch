// Package notify implements webhook notifications for detected configuration
// drift. When driftwatch identifies differences between a deployed service and
// its source definition it can optionally POST a structured JSON payload to a
// user-supplied webhook URL, enabling integration with alerting systems such as
// Slack, PagerDuty, or custom HTTP endpoints.
//
// Basic usage:
//
//	n := notify.New("https://hooks.example.com/drift")
//	err := n.Send(notify.Payload{
//		Service:    "api-gateway",
//		DriftCount: 2,
//		Fields:     []string{"replicas", "image"},
//	})
package notify
