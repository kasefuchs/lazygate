package nomad

import (
	"github.com/go-logr/logr"
	"github.com/hashicorp/nomad/api"
	"github.com/kasefuchs/lazygate/pkg/provider"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var _ provider.Provider = (*Provider)(nil)

// Provider based on Nomad API.
type Provider struct {
	log    logr.Logger // Provider logger.
	client *api.Client // Nomad API client.
}

func (p *Provider) initClient() error {
	var err error

	cfg := api.DefaultConfig()
	p.client, err = api.NewClient(cfg)

	return err
}

func (p *Provider) Init(opt *provider.InitOptions) error {
	p.log = opt.Log

	if err := p.initClient(); err != nil {
		return err
	}

	p.log.Info("initialized provider")
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
