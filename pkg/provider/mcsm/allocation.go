package mcsm

import (
	"fmt"

	"github.com/kasefuchs/lazygate/pkg/provider"
)

var _ provider.Allocation = (*Allocation)(nil)

// Allocation represents Docker provider item.
type Allocation struct {
	item           *item
	client         *Client
	configFilePath string
}

func NewAllocation(client *Client, item *item) *Allocation {
	return &Allocation{
		item:   item,
		client: client,
	}
}

func (a *Allocation) Stop() error {
	return a.client.ServerStop(a.item.instance.Id, a.client.daemon_id)
}

func (a *Allocation) Start() error {
	return a.client.ServerStart(a.item.instance.Id, a.client.daemon_id)
}

func (a *Allocation) State() provider.AllocationState {
	inst, err := a.client.GetInstance(a.item.instance.Id, a.client.daemon_id)
	if err != nil {
		fmt.Errorf("failed to get isntance for state %v", err)
	}
	switch inst.Data.Status {
	case 0:
	case 1:
		return provider.AllocationStateStopped
	case 2:
	case 3:
		return provider.AllocationStateStarted
	case -1:
		return provider.AllocationStateUnknown
	}
	return provider.AllocationStateUnknown
}

func (a *Allocation) ParseConfig(cfg interface{}, rootLabel string) (interface{}, error) {
	inst, err := a.client.GetInstance(a.item.instance.Id, a.client.daemon_id)
	if err != nil {
		return nil, err
	}

	cfgtype := cfg.(*provider.AllocationConfig)
	cfgtype.Server = inst.Data.Config.Nickname
	// for _, tag := range inst.Data.Config.Tag {
	// 	if strings.HasPrefix(tag, "lazy.name=") {
	// 		cfg2.Server = tag[10 : len(tag)-1]
	// 		break
	// 	}
	// }
	// fill the rest with default options

	return cfgtype, nil
}
