package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/kasefuchs/lazygate/pkg/config/allocation"
)

// Allocation internal data in Docker context.
type item struct {
	container *types.Container // Server container.
}

func (p *Provider) itemList() ([]*item, error) {
	var items []*item

	containerList, err := p.client.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	for _, cnt := range containerList {
		if _, err := allocation.ParseLabels(cnt.Labels); err != nil {
			continue
		}

		items = append(items, &item{
			container: &cnt,
		})
	}

	return items, nil
}
