package mcsm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pufferpanel/pufferpanel/v3"
)

type HeaderTransport struct {
	Base    http.RoundTripper
	Headers map[string]string
}

func (h *HeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Clone request to avoid modifying the original
	newReq := req.Clone(req.Context())

	// Inject default headers (only if not already set)
	for key, value := range h.Headers {
		if newReq.Header.Get(key) == "" {
			newReq.Header.Set(key, value)
		}
	}
	return h.Base.RoundTrip(newReq)
}

type Client struct {
	ctx        context.Context
	api_client *http.Client
	daemon_id  string
	base_url   string
}

func NewClient(ctx context.Context, baseUrl, api_key string, daemon_id string) *Client {
	api_client := &http.Client{
		Transport: &HeaderTransport{
			Headers: map[string]string{
				"apikey":           api_key,
				"Content-Type":     "application/json; charset=utf-8",
				"X-Requested-With": "XMLHttpRequest",
			},
		},
	}
	return &Client{
		ctx:        ctx,
		api_client: api_client,
		daemon_id:  daemon_id,
		base_url:   baseUrl,
	}
}

func (c *Client) ServerStop(id string) error {
	req, err := http.NewRequestWithContext(c.ctx, http.MethodPost, fmt.Sprintf("%s/api/servers/%s/stop", c.base_url, id), nil)
	if err != nil {
		return fmt.Errorf("creating server stop request: %w", err)
	}

	resp, err := c.api_client.Do(req)
	if err != nil {
		return fmt.Errorf("executing server stop request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	return nil
}

func (c *Client) ServerStart(id string) error {

	req, err := http.NewRequestWithContext(c.ctx, http.MethodPost, fmt.Sprintf("%s/api/servers/%s/start", c.base_url, id), nil)
	if err != nil {
		return fmt.Errorf("creating server start request: %w", err)
	}

	resp, err := c.api_client.Do(req)
	if err != nil {
		return fmt.Errorf("executing server start request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	return nil
}

func (c *Client) ServerStatus(id string) (*pufferpanel.ServerRunning, error) {

	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, fmt.Sprintf("%s/api/servers/%s/status", c.base_url, id), nil)
	if err != nil {
		return nil, fmt.Errorf("creating server status request: %w", err)
	}

	resp, err := c.api_client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing server status request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var status *pufferpanel.ServerRunning
	if err := json.Unmarshal(body, &status); err != nil {
		return nil, err
	}

	return status, nil
}

func (c *Client) ServerSearch() (*InstanceSearch, error) {
	return c.serverSearch(1, 50)
}

func (c *Client) serverSearch(page uint, page_size uint) (*InstanceSearch, error) {
	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, fmt.Sprintf("%s/api/service/remote_service_instances?daemonId=%s&page=%d&page_size=%d", c.base_url, c.daemon_id, page, page_size), nil)
	if err != nil {
		return nil, fmt.Errorf("creating server search request: %w", err)
	}

	resp, err := c.api_client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing server search request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var search *InstanceSearch
	if err := json.Unmarshal(body, &search); err != nil {
		return nil, err
	}
	if uint(search.Data.Page) < uint(search.Data.MaxPage) {
		extraSearch, err := c.serverSearch(page+1, 50)
		if err != nil {
			return nil, fmt.Errorf("server search page %d: %w", page+1, err)
		}
		search.Data.Data = append(search.Data.Data, extraSearch.Data.Data...)
	}

	return search, nil
}

func (c *Client) GetInstance(id, daemon_id string) (*InstanceDetail, error) {
	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, fmt.Sprintf("%s/api/instance?uuid=%s&daemonId=%s", c.base_url, id, daemon_id), nil)
	if err != nil {
		return nil, fmt.Errorf("creating get instance request: %w", err)
	}

	resp, err := c.api_client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing get instance request: %w", err)
	}
	defer resp.Body.Close()
	var res InstanceDetail
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, fmt.Errorf("Failed to Decode Instace List: %w", err)
	}
	return &res, nil
}
