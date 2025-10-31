package pufferpanel

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/pufferpanel/pufferpanel/v3"
	"github.com/pufferpanel/pufferpanel/v3/models"
	"github.com/pufferpanel/pufferpanel/v3/oauth2"
)

type Client struct {
	ctx            context.Context
	baseUrl        string
	clientId       string
	clientSecret   string
	token          *oauth2.TokenResponse
	tokenExpiresAt time.Time
}

func NewClient(ctx context.Context, baseUrl, clientId, clientSecret string) *Client {
	return &Client{
		ctx:          ctx,
		baseUrl:      baseUrl,
		clientId:     clientId,
		clientSecret: clientSecret,
	}
}

func (c *Client) RequestToken() (*oauth2.TokenResponse, error) {
	form := url.Values{
		"client_id":     {c.clientId},
		"client_secret": {c.clientSecret},
		"grant_type":    {"client_credentials"},
	}

	req, err := http.NewRequestWithContext(c.ctx, http.MethodPost, fmt.Sprintf("%s/oauth2/token", c.baseUrl), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating token request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing token request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	var token oauth2.TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, fmt.Errorf("decoding token response: %w", err)
	}

	return &token, nil
}

func (c *Client) RenewToken() (*oauth2.TokenResponse, error) {
	if c.token != nil && time.Now().Before(c.tokenExpiresAt) {
		return c.token, nil
	}

	tok, err := c.RequestToken()
	if err != nil {
		return nil, err
	}

	c.token = tok
	c.tokenExpiresAt = time.Now().Add(time.Second * time.Duration(tok.ExpiresIn))

	return c.token, nil
}

func (c *Client) ServerStop(id string) error {
	token, err := c.RenewToken()
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(c.ctx, http.MethodPost, fmt.Sprintf("%s/api/servers/%s/stop", c.baseUrl, id), nil)
	if err != nil {
		return fmt.Errorf("creating server stop request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing server stop request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	return nil
}

func (c *Client) ServerStart(id string) error {
	token, err := c.RenewToken()
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(c.ctx, http.MethodPost, fmt.Sprintf("%s/api/servers/%s/start", c.baseUrl, id), nil)
	if err != nil {
		return fmt.Errorf("creating server start request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing server start request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	return nil
}

func (c *Client) ServerStatus(id string) (*pufferpanel.ServerRunning, error) {
	token, err := c.RenewToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, fmt.Sprintf("%s/api/servers/%s/status", c.baseUrl, id), nil)
	if err != nil {
		return nil, fmt.Errorf("creating server status request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))

	resp, err := http.DefaultClient.Do(req)
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

func (c *Client) ServerSearch() (*models.ServerSearchResponse, error) {
	return c.serverSearch(1)
}

func (c *Client) serverSearch(page uint) (*models.ServerSearchResponse, error) {
	token, err := c.RenewToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, fmt.Sprintf("%s/api/servers?page=%d&limit=100", c.baseUrl, page), nil)
	if err != nil {
		return nil, fmt.Errorf("creating server search request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing server search request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var search *models.ServerSearchResponse
	if err := json.Unmarshal(body, &search); err != nil {
		return nil, err
	}
	if int64(search.Metadata.Paging.Size*page) < search.Metadata.Paging.Total {
		var extraSearch *models.ServerSearchResponse
		extraSearch, err = c.serverSearch(page + 1)
		if err != nil {
			return nil, fmt.Errorf("server search page %d: %w", page+1, err)
		}
		search.Servers = append(search.Servers, extraSearch.Servers...)
	}

	return search, nil
}

func (c *Client) ReadServerFile(id, path string) ([]byte, error) {
	token, err := c.RenewToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(c.ctx, http.MethodGet, fmt.Sprintf("%s/api/servers/%s/file/%s", c.baseUrl, id, path), nil)
	if err != nil {
		return nil, fmt.Errorf("creating read server file request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("%s %s", token.TokenType, token.AccessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing read server file request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	return io.ReadAll(resp.Body)
}
