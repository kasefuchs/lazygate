package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/kasefuchs/lazygate/pkg/config/dynamic"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

// Allocation internal data in Docker context.
type item struct {
	config    *dynamic.Config  // Server dynamic configuration.
	container *types.Container // Server container.
}

func (p *Provider) itemList() ([]*item, error) {
	var items []*item

	containerList, err := p.client.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	for _, cnt := range containerList {
		cfg, err := dynamic.ParseLabels(cnt.Labels)
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

	for _, i := range items {
		if i.config.Server == srv.ServerInfo().Name() {
			return i, nil
		}
	}

	return nil, fmt.Errorf("container not found")
}
