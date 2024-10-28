package nomad

import (
	"github.com/hashicorp/nomad/api"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

var (
	countUp   int64 = 1
	countDown int64 = 0
)

type Provider struct {
	client *api.Client
}

func (p *Provider) Init() error {
	var err error

	cfg := api.DefaultConfig()
	p.client, err = api.NewClient(cfg)

	return err
}

func (p *Provider) Pause(srv *proxy.RegisteredServer) error {
	return p.scale(srv, &countDown)
}

func (p *Provider) Resume(srv *proxy.RegisteredServer) error {
	return p.scale(srv, &countUp)
}
