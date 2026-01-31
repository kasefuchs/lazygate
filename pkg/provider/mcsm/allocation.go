package mcsm

import (
	"github.com/kasefuchs/lazygate/pkg/provider"
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
	return a.client.ServerStop(a.item.instance.Id)
}

func (a *Allocation) Start() error {
	return a.client.ServerStart(a.item.instance.Id)
}

func (a *Allocation) State() provider.AllocationState {
	status, err := a.client.ServerStatus(a.item.instance.Id)
	if err != nil {
		return provider.AllocationStateUnknown
	}

	if status.Running {
		return provider.AllocationStateStarted
	}

	return provider.AllocationStateStopped
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
