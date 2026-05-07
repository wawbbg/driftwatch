// Package resolve maps service names to their live config endpoints
// and expected config definitions, combining fetcher and config sources.
package resolve

import (
	"context"
	"fmt"

	"github.com/driftwatch/internal/config"
	"github.com/driftwatch/internal/fetcher"
)

// Result holds the resolved live and expected configs for a service.
type Result struct {
	ServiceName string
	Live        map[string]string
	Expected    map[string]string
}

// Resolver resolves service configs.
type Resolver struct {
	fetcher *fetcher.Fetcher
}

// New returns a Resolver using the default HTTP fetcher.
func New() *Resolver {
	return &Resolver{fetcher: fetcher.New()}
}

// NewWithFetcher returns a Resolver using the provided fetcher.
func NewWithFetcher(f *fetcher.Fetcher) *Resolver {
	return &Resolver{fetcher: f}
}

// Service resolves a single service by fetching its live config and
// pairing it with the expected fields from the config definition.
func (r *Resolver) Service(ctx context.Context, svc config.Service) (Result, error) {
	if svc.URL == "" {
		return Result{}, fmt.Errorf("resolve: service %q has no URL", svc.Name)
	}

	live, err := r.fetcher.Fetch(ctx, svc.URL)
	if err != nil {
		return Result{}, fmt.Errorf("resolve: fetch %q: %w", svc.Name, err)
	}

	return Result{
		ServiceName: svc.Name,
		Live:        live,
		Expected:    svc.Expected,
	}, nil
}

// All resolves every service in the provided config, collecting errors
// without aborting early. Callers receive all results and all errors.
func (r *Resolver) All(ctx context.Context, cfg *config.Config) ([]Result, []error) {
	var results []Result
	var errs []error

	for _, svc := range cfg.Services {
		res, err := r.Service(ctx, svc)
		if err != nil {
			errs = append(errs, err)
			continue
		}
		results = append(results, res)
	}

	return results, errs
}
