package pufferpanel

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/kasefuchs/lazygate/pkg/provider"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var _ provider.Provider = (*Provider)(nil)

const name = "pufferpanel"

// Provider based on Docker API.
type Provider struct {
	log    logr.Logger     // Provider logger.
	ctx    context.Context // Provider context.
	client *Client         // PufferPanel API client.
	config *Config         // Provider config.
}

func (p *Provider) Name() string {
	return name
}

func (p *Provider) DefaultConfig() interface{} {
	return &Config{
		ConfigFilePath: "lazygate.json",
	}
}

func (p *Provider) initClient() {
	p.client = NewClient(p.ctx, p.config.BaseUrl, p.config.ClientId, p.config.ClientSecret)
}

func (p *Provider) Init(opt *provider.InitOptions) error {
	p.log = logr.FromContextOrDiscard(opt.Ctx).WithName(provider.LogName)
	p.ctx = opt.Ctx
	p.config = opt.Config.(*Config)

	p.initClient()

	p.log.Info("initialized")
	return nil
}

func (p *Provider) AllocationGet(srv proxy.RegisteredServer) (provider.Allocation, error) {
	allocs, err := p.AllocationList()
	if err != nil {
		return nil, err
	}

	for _, alloc := range allocs {
		cfg, err := provider.ParseAllocationConfig(alloc)
		if err != nil {
			continue
		}

		if cfg.Server == srv.ServerInfo().Name() {
			return alloc, nil
		}
	}

	return nil, fmt.Errorf("no allocation found")
}

func (p *Provider) AllocationList() ([]provider.Allocation, error) {
	items, err := p.itemList()
	if err != nil {
		return nil, err
	}

	allocs := make([]provider.Allocation, len(items))
	for i, it := range items {
		allocs[i] = NewAllocation(p.client, it, p.config.ConfigFilePath)
	}

	return allocs, nil
}
