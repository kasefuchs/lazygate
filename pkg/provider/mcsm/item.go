package mcsm

import "slices"

// Allocation internal data in PufferPanel context.
type instance struct {
	Id string
}
type item struct {
	instance *instance
}

func (p *Provider) itemList() ([]*item, error) {
	search, err := p.client.ServerSearch()
	if err != nil {
		return nil, err
	}

	var items []*item
	for _, i := range search.Data.Data {
		if (slices.Contains(i.Config.Tag, "lazy")) == false {
			continue
		}
		items = append(items, &item{
			instance: &instance{Id: i.InstanceUUID},
		})
	}

	return items, nil
}
