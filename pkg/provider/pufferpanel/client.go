package pufferpanel

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type Client struct {
	baseUrl       string
	client_id     string
	client_secret string
}

type OAuth2Token struct {
	Access_token string  `json:"access_token"`
	Expires_in   float64 `json:"expires_in"`
	Scope        string  `json:"scope"`
	Token_type   string  `json:"token_type"`
}

func NewClient() (*Client, error) {
	baseUrl := os.Getenv("LAZYGATE_PUFFERPANEL_URL")
	client_id := os.Getenv("LAZYGATE_PUFFERPANEL_CLIENTID")
	client_secret := os.Getenv("LAZYGATE_PUFFERPANEL_CLIENTSECRET")
	return &Client{
		baseUrl:       strings.TrimSpace(baseUrl),
		client_id:     strings.TrimSpace(client_id),
		client_secret: strings.TrimSpace(client_secret),
	}, nil
}

func (client *Client) getToken() (*OAuth2Token, error) {
	var token OAuth2Token
	data := url.Values{}

	data.Set("client_id", client.client_id)
	data.Set("client_secret", client.client_secret)
	data.Set("grant_type", "client_credentials")
	resp, err := http.PostForm(client.baseUrl+"oauth2/token", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}
