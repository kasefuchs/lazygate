package mcsm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HeaderTransport struct {
	Headers map[string]string
	Apikey  string
}

func (h *HeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range h.Headers {
		if req.Header.Get(key) == "" {
			req.Header[key] = []string{value}
		}
	}
	q := req.URL.Query()
	q.Add("apikey", h.Apikey)
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		defer res.Body.Close()
		resp, errp := io.ReadAll(res.Body)
		if errp != nil {
			return nil, fmt.Errorf("Error reading io: %s", errp)
		}
		return nil, fmt.Errorf("%s\n%s", res.Status, string(resp))
	}
	return res, nil
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
				"Content-Type":     "application/json; charset=utf-8",
				"X-Requested-With": "XMLHttpRequest",
			},
			Apikey: api_key,
		},
	}
	return &Client{
		ctx:        ctx,
		api_client: api_client,
		daemon_id:  daemon_id,
		base_url:   baseUrl,
	}
}

func (c *Client) ServerStop(id string, did string) error {
	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, fmt.Sprintf("%s/api/protected_instance/stop?uuid=%s&daemonId=%s", c.base_url, id, did), nil)
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

func (c *Client) ServerStart(id string, did string) error {

	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, fmt.Sprintf("%s/api/protected_instance/open?uuid=%s&daemonId=%s", c.base_url, id, did), nil)
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
