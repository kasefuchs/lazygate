package docker

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

type Provider struct {
	client *client.Client
}

func (p *Provider) Init() error {
	var err error

	p.client, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	return nil
}

func (p *Provider) Pause(srv *proxy.RegisteredServer) error {
	cnt, err := p.getContainer(srv)
	if err != nil {
		return err
	}

	return p.client.ContainerStop(context.Background(), cnt.ID, container.StopOptions{})
}

func (p *Provider) Resume(srv *proxy.RegisteredServer) error {
	cnt, err := p.getContainer(srv)
	if err != nil {
		return err
	}

	return p.client.ContainerStart(context.Background(), cnt.ID, container.StartOptions{})
}
