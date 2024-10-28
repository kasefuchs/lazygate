package nomad

import (
	"fmt"

	"github.com/hashicorp/nomad/api"
	"github.com/kasefuchs/lazygate/pkg/config/dynamic"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

type alloc struct {
	job       *api.Job
	taskGroup *api.TaskGroup
}

func (p *Provider) getAlloc(srv *proxy.RegisteredServer) (*alloc, error) {
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
				cfg, err := dynamic.ParseTags(service.Tags)
				if err != nil || cfg.Server == "" {
					continue
				}

				if srv == nil || cfg.Server == (*srv).ServerInfo().Name() {
					return &alloc{
						job:       job,
						taskGroup: taskGroup,
					}, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("alloc not found")
}

func (p *Provider) scale(srv *proxy.RegisteredServer, count *int64) error {
	alloc, err := p.getAlloc(srv)
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
