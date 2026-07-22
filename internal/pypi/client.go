package pypi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"forgejo.eng.ultra-cis.com/devops/packmule/internal/consts"
)

type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// IndexFetcher fetches the full list of projects from the simple index.
type IndexFetcher interface {
	FetchIndex(ctx context.Context) (*Index, error)
}

// ProjectFetcher fetches file/version details for a single project.
type ProjectFetcher interface {
	FetchProject(ctx context.Context, name string) (*ProjectDetail, error)
}

// Store persists collected data. Implementations might write to a
// database, a file, an object store, etc.
type Store interface {
	SaveIndex(ctx context.Context, idx *Index) error
	SaveProject(ctx context.Context, detail *ProjectDetail) error
}

type Client struct {
	HTTP      HTTPDoer
	BaseURL   string // defaults to https://pypi.org/simple/
	UserAgent string
}

// NewClient builds a Client with sane defaults (30s timeout).
func NewClient() *Client {
	return &Client{
		HTTP:      &http.Client{Timeout: 30 * time.Second},
		BaseURL:   consts.SimpleIndexURL,
		UserAgent: consts.DefaultUserAgent,
	}
}

func (c *Client) newRequest(ctx context.Context, url string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", consts.SimpleJSONAccept)
	req.Header.Set("User-Agent", c.UserAgent)
	return req, nil
}

// FetchIndex retrieves and decodes the full project list.
func (c *Client) FetchIndex(ctx context.Context) (*Index, error) {
	req, err := c.newRequest(ctx, c.baseURL())
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching index: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<10))
		return nil, fmt.Errorf("unexpected status %d fetching index: %s", resp.StatusCode, body)
	}

	var idx Index
	if err := json.NewDecoder(resp.Body).Decode(&idx); err != nil {
		return nil, fmt.Errorf("decoding index: %w", err)
	}
	return &idx, nil
}

// FetchProject retrieves file/version details for a single project.
func (c *Client) FetchProject(ctx context.Context, name string) (*ProjectDetail, error) {
	url := c.baseURL() + name + "/"
	req, err := c.newRequest(ctx, url)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching project %q: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1<<10))
		return nil, fmt.Errorf("unexpected status %d fetching project %q: %s", resp.StatusCode, name, body)
	}

	var detail ProjectDetail
	if err := json.NewDecoder(resp.Body).Decode(&detail); err != nil {
		return nil, fmt.Errorf("decoding project %q: %w", name, err)
	}
	return &detail, nil
}

func (c *Client) baseURL() string {
	if c.BaseURL != "" {
		return c.BaseURL
	}
	return consts.SimpleIndexURL
}
