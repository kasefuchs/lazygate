package pufferpanel

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/kasefuchs/lazygate/pkg/config/allocation"
	"github.com/kasefuchs/lazygate/pkg/provider"
)

var _ provider.Allocation = (*Allocation)(nil)

// Allocation represents Docker provider item.
type Allocation struct {
	client *Client
	item   *item
}

func NewAllocation(client *Client, item *item) *Allocation {
	return &Allocation{
		client: client,
		item:   item,
	}
}

func (a *Allocation) Stop() error {
	token, err := a.client.getToken()
	if err != nil {
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", a.client.baseUrl+"api/servers/"+a.item.id+"/stop", nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", token.Token_type+" "+token.Access_token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (a *Allocation) Start() error {
	token, err := a.client.getToken()
	if err != nil {
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", a.client.baseUrl+"api/servers/"+a.item.id+"/start", nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", token.Token_type+" "+token.Access_token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (a *Allocation) State() provider.AllocationState {
	type Status struct {
		Running    bool `json:"running"`
		Installing bool `json:"installing"`
	}
	token, err := a.client.getToken()
	if err != nil {
		return provider.AllocationStateUnknown
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", a.client.baseUrl+"api/servers/"+a.item.id+"/stop", nil)
	if err != nil {
		return provider.AllocationStateUnknown
	}
	req.Header.Add("Authorization", token.Token_type+" "+token.Access_token)
	resp, err := client.Do(req)
	if err != nil {
		return provider.AllocationStateUnknown
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return provider.AllocationStateUnknown
	}
	var status Status
	json.Unmarshal(body, &status)
	if status.Running {
		return provider.AllocationStateStarted
	}

	return provider.AllocationStateStopped
}

func (a *Allocation) Config() (*allocation.Config, error) {
	config := map[string]string{}
	token, err := a.client.getToken()
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", a.client.baseUrl+"api/servers/"+a.item.id+"/file/lazygate.json", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", token.Token_type+" "+token.Access_token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	json.Unmarshal(body, &config)

	//fmt.Println(inspect.Config.Labels)
	return allocation.ParseLabels(config)
}
