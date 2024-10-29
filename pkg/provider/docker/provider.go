package docker

import (
	"github.com/docker/docker/client"
	"github.com/kasefuchs/lazygate/pkg/provider"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var _ provider.Provider = (*Provider)(nil)

// Provider based on Docker API.
type Provider struct {
	client *client.Client // Docker API client.
}

func (p *Provider) Init(opt provider.InitOptions) error {
	var err error

	p.client, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) AllocationGet(srv proxy.RegisteredServer) (provider.Allocation, error) {
	it, err := p.itemGet(srv)
	if err != nil {
		return nil, err
	}

	return NewAllocation(p.client, it), nil
}

func (p *Provider) AllocationList() ([]provider.Allocation, error) {
	items, err := p.itemList()
	if err != nil {
		return nil, err
	}

	allocs := make([]provider.Allocation, len(items))
	for i, it := range items {
		allocs[i] = NewAllocation(p.client, it)
	}

	return allocs, nil
}