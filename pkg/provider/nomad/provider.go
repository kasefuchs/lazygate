package nomad

import (
	"fmt"

	"github.com/go-logr/logr"
	"github.com/hashicorp/nomad/api"
	"github.com/kasefuchs/lazygate/pkg/provider"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var _ provider.Provider = (*Provider)(nil)

const name = "nomad"

// Provider based on Nomad API.
type Provider struct {
	log    logr.Logger // Provider logger.
	client *api.Client // Nomad API client.
	config *Config     // Provider config.
}

func (p *Provider) initClient() error {
	var err error

	cfg := api.DefaultConfig()
	p.client, err = api.NewClient(cfg)

	return err
}

func (p *Provider) Init(opt *provider.InitOptions) error {
	p.log = logr.FromContextOrDiscard(opt.Ctx).WithName(provider.LogName)
	p.config = opt.Config.(*Config)

	if err := p.initClient(); err != nil {
		return err
	}

	p.log.Info("initialized")
	return nil
}

func (p *Provider) Name() string {
	return name
}

func (p *Provider) DefaultConfig() interface{} {
	return &Config{}
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
		allocs[i] = NewAllocation(p.client, it)
	}

	return allocs, nil
}
