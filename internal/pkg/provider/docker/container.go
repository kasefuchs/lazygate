package docker

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/kasefuchs/lazygate/internal/pkg/config/dynamic"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

func (p *Provider) getContainer(srv *proxy.RegisteredServer) (*types.Container, error) {
	containerList, err := p.client.ContainerList(context.Background(), container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	for _, cnt := range containerList {
		cfg, err := dynamic.ParseLabels(cnt.Labels)
		if err != nil || cfg.Server == "" {
			continue
		}

		if srv == nil || cfg.Server == (*srv).ServerInfo().Name() {
			return &cnt, nil
		}
	}

	return nil, fmt.Errorf("container not found")
}
