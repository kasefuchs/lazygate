package nomad

import (
	"fmt"
	"github.com/hashicorp/nomad/api"
	"github.com/kasefuchs/lazygate/internal/pkg/config"
	"github.com/kasefuchs/lazygate/internal/pkg/provider"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

const providerName = "nomad"

var _ provider.Provider = (*Provider)(nil)

var (
	countUp   int64 = 1
	countDown int64 = 0
)

type allocation struct {
	job       *api.Job
	taskGroup *api.TaskGroup
}

type Provider struct {
	client *api.Client
}

func (p *Provider) Name() string {
	return providerName
}

func (p *Provider) Init() error {
	var err error

	cfg := api.DefaultConfig()
	p.client, err = api.NewClient(cfg)

	return err
}

func (p *Provider) getAllocation(srv *proxy.RegisteredServer) (*allocation, error) {
	jobStubs, _, err := p.client.Jobs().ListOptions(nil, nil)
	if err != nil {
		return nil, err
	}

	for _, jobStub := range jobStubs {
		job, _, err := p.client.Jobs().Info(jobStub.ID, nil)
		if err != nil {
			return nil, err
		}

		for _, taskGroup := range job.TaskGroups {
			for _, service := range taskGroup.Services {
				cfg, err := config.ParseTags(service.Tags)
				if err != nil || cfg.Server.Name == "" {
					continue
				}

				if srv == nil || cfg.Server.Name == (*srv).ServerInfo().Name() {
					return &allocation{
						job:       job,
						taskGroup: taskGroup,
					}, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("allocation not found")
}

func (p *Provider) scale(srv *proxy.RegisteredServer, count *int64) error {
	alloc, err := p.getAllocation(srv)
	if err != nil {
		return err
	}

	if *count != int64(*alloc.taskGroup.Count) {
		req := &api.ScalingRequest{
			Count: count,
			Target: map[string]string{
				"Job":   *alloc.job.ID,
				"Group": *alloc.taskGroup.Name,
			},
			PolicyOverride: false,
			JobModifyIndex: 0,
		}
		if _, _, err := p.client.Jobs().ScaleWithRequest(*alloc.job.ID, req, nil); err != nil {
			return err
		}
	}

	return nil
}

func (p *Provider) Pause(srv proxy.RegisteredServer) error {
	return p.scale(&srv, &countDown)
}

func (p *Provider) Resume(srv proxy.RegisteredServer) error {
	return p.scale(&srv, &countUp)
}

func (p *Provider) ResumeAny() error {
	return p.scale(nil, &countUp)
}
