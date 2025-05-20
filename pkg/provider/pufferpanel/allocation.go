package pufferpanel

import (
	"encoding/json"

	"github.com/kasefuchs/lazygate/pkg/provider"
	"github.com/kasefuchs/lazygate/pkg/utils"
)

var _ provider.Allocation = (*Allocation)(nil)

// Allocation represents Docker provider item.
type Allocation struct {
	item           *item
	client         *Client
	configFilePath string
}

func NewAllocation(client *Client, item *item, configFilePath string) *Allocation {
	return &Allocation{
		item:           item,
		client:         client,
		configFilePath: configFilePath,
	}
}

func (a *Allocation) Stop() error {
	return a.client.ServerStop(a.item.server.Identifier)
}

func (a *Allocation) Start() error {
	return a.client.ServerStart(a.item.server.Identifier)
}

func (a *Allocation) State() provider.AllocationState {
	status, err := a.client.ServerStatus(a.item.server.Identifier)
	if err != nil {
		return provider.AllocationStateUnknown
	}

	if status.Running {
		return provider.AllocationStateStarted
	}

	return provider.AllocationStateStopped
}

func (a *Allocation) ParseConfig(cfg interface{}, rootLabel string) (interface{}, error) {
	bytes, err := a.client.ReadServerFile(a.item.server.Identifier, a.configFilePath)
	if err != nil {
		return nil, err
	}

	labels := map[string]string{}
	if err := json.Unmarshal(bytes, &labels); err != nil {
		return nil, err
	}

	return utils.ParseLabels(labels, cfg, rootLabel)
}
