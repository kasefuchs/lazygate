package pufferpanel

import "github.com/pufferpanel/pufferpanel/v3/models"

// Allocation internal data in PufferPanel context.
type item struct {
	server *models.ServerView
}

func (p *Provider) itemList() ([]*item, error) {
	search, err := p.client.ServerSearch()
	if err != nil {
		return nil, err
	}

	var items []*item
	for _, server := range search.Servers {
		items = append(items, &item{
			server: server,
		})
	}

	return items, nil
}
