package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/kasefuchs/lazygate/pkg/provider"
	"github.com/kasefuchs/lazygate/pkg/utils"
)

var _ provider.Allocation = (*Allocation)(nil)

// Allocation represents Docker provider item.
type Allocation struct {
	client *client.Client // Docker API client.
	item   *item          // Docker container item.
}

func NewAllocation(client *client.Client, item *item) *Allocation {
	return &Allocation{
		client: client,
		item:   item,
	}
}

func (a *Allocation) Stop() error {
	return a.client.ContainerStop(context.Background(), a.item.container.ID, container.StopOptions{})
}

func (a *Allocation) Start() error {
	return a.client.ContainerStart(context.Background(), a.item.container.ID, container.StartOptions{})
}

func (a *Allocation) State() provider.AllocationState {
	inspect, err := a.client.ContainerInspect(context.Background(), a.item.container.ID)
	if err != nil {
		return provider.AllocationStateUnknown
	}

	if inspect.State.Running {
		return provider.AllocationStateStarted
	}

	return provider.AllocationStateStopped
}

func (a *Allocation) ParseConfig(cfg interface{}, rootLabel string) (interface{}, error) {
	inspect, err := a.client.ContainerInspect(context.Background(), a.item.container.ID)
	if err != nil {
		return nil, err
	}

	return utils.ParseLabels(inspect.Config.Labels, cfg, rootLabel)
}
