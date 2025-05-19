package pufferpanel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Allocation internal data in Docker context.
type item struct {
	id string
}
type Servers struct {
	Servers []map[string]interface{} `json:"servers"`
	Paging  map[string]interface{}   `json:"paging"`
}

func (p *Provider) itemList() ([]*item, error) {
	var items []*item

	token, err := p.client.getToken()
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	var req *http.Request
	req, err = http.NewRequest("GET", p.client.baseUrl+"api/servers", nil)
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
	var servers Servers
	err = json.Unmarshal(body, &servers)
	if err != nil {
		return nil, err
	}
	for _, server := range servers.Servers {
		if id := server["id"]; id != nil {
			item := &item{id: id.(string)}
			fmt.Println(item.id)
			items = append(items, item) //maybe not the right cast
		}
	}

	return items, nil
}
