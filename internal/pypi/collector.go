package pypi

import (
	"context"
	"fmt"
	"sync"
)

type CollectorConfig struct {
	// Concurrency is the number of simultaneous per-project detail
	// fetches. Ignored if FetchDetails is false.
	Concurrency int
	// FetchDetails, if true, also fetches file/version details for
	// every project in the index (this is a lot of requests -- PyPI
	// hosts several hundred thousand projects).
	FetchDetails bool
}

// Collector ties an IndexFetcher/ProjectFetcher pair to a Store and runs
// a full (or index-only) collection pass.
type Collector struct {
	Index   IndexFetcher
	Project ProjectFetcher
	Store   Store
	Config  CollectorConfig
}

// NewCollector builds a Collector using the default HTTP Client for both
// index and project fetching.
func NewCollector(store Store, cfg CollectorConfig) *Collector {
	client := NewClient()
	if cfg.Concurrency <= 0 {
		cfg.Concurrency = 10
	}
	return &Collector{
		Index:   client,
		Project: client,
		Store:   store,
		Config:  cfg,
	}
}

// Run fetches the index, saves it, and (if configured) fetches and saves
// per-project details concurrently. It stops at the first index-fetch
// error but collects and returns per-project errors without aborting the
// whole run.
func (c *Collector) Run(ctx context.Context) error {
	idx, err := c.Index.FetchIndex(ctx)
	if err != nil {
		return fmt.Errorf("collector: fetching index: %w", err)
	}
	if err := c.Store.SaveIndex(ctx, idx); err != nil {
		return fmt.Errorf("collector: saving index: %w", err)
	}

	if !c.Config.FetchDetails {
		return nil
	}

	sem := make(chan struct{}, c.Config.Concurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errs []error

	for _, p := range idx.Projects {
		p := p
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			detail, err := c.Project.FetchProject(ctx, p.Name)
			if err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("project %q: %w", p.Name, err))
				mu.Unlock()
				return
			}
			if err := c.Store.SaveProject(ctx, detail); err != nil {
				mu.Lock()
				errs = append(errs, fmt.Errorf("project %q: save: %w", p.Name, err))
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("collector: %d project errors, e.g. %w", len(errs), errs[0])
	}
	return nil
}
