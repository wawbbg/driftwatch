// Package filter provides service filtering utilities for driftwatch.
package filter

import "strings"

// Options holds filtering criteria for services.
type Options struct {
	Names  []string
	Labels map[string]string
}

// Service represents a minimal view of a service for filtering.
type Service struct {
	Name   string
	Labels map[string]string
}

// Match reports whether the service matches the given filter options.
// An empty Options matches all services.
func Match(svc Service, opts Options) bool {
	if len(opts.Names) > 0 && !containsName(opts.Names, svc.Name) {
		return false
	}
	for k, v := range opts.Labels {
		got, ok := svc.Labels[k]
		if !ok || !strings.EqualFold(got, v) {
			return false
		}
	}
	return true
}

// Apply filters a slice of services using the given options.
func Apply(services []Service, opts Options) []Service {
	if len(opts.Names) == 0 && len(opts.Labels) == 0 {
		return services
	}
	out := make([]Service, 0, len(services))
	for _, svc := range services {
		if Match(svc, opts) {
			out = append(out, svc)
		}
	}
	return out
}

func containsName(names []string, name string) bool {
	for _, n := range names {
		if strings.EqualFold(n, name) {
			return true
		}
	}
	return false
}
