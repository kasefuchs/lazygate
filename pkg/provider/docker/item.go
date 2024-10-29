package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/kasefuchs/lazygate/pkg/config"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

// Allocation internal data in Docker context.
type item struct {
	config    *config.Config   // Server dynamic configuration.
	container *types.Container // Server container.
}

func (p *Provider) itemList() ([]*item, error) {
	var items []*item

	containerList, err := p.client.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	for _, cnt := range containerList {
		cfg, err := config.ParseLabels(cnt.Labels)
		if err != nil {
			continue
		}

		items = append(items, &item{
			config:    cfg,
			container: &cnt,
		})
	}

	return items, nil
}

func (p *Provider) itemGet(srv proxy.RegisteredServer) (*item, error) {
	items, err := p.itemList()
	if err != nil {
		return nil, err
	}

	for _, it := range items {
		if it.config.Server == srv.ServerInfo().Name() {
			return it, nil
		}
	}

	return nil, fmt.Errorf("item not found")
}
