package nomad

import (
	"fmt"

	"github.com/hashicorp/nomad/api"
	"github.com/kasefuchs/lazygate/pkg/config"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

// Allocation internal data in Nomad context.
type item struct {
	job     *api.Job       // Task job
	group   *api.TaskGroup // Task group
	config  *config.Config // Server dynamic configuration.
	service *api.Service   // Task service.
}

func (p *Provider) itemList() ([]*item, error) {
	var items []*item

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
				if err != nil {
					continue
				}

				items = append(items, &item{
					job:     job,
					group:   taskGroup,
					config:  cfg,
					service: service,
				})
			}
		}
	}

	return items, nil
}

func (p *Provider) itemGet(srv proxy.RegisteredServer) (*item, error) {
	items, err := p.itemList()
	if err != nil {
		return nil, err
	}

	for _, it := range items {
		if it.config.Server == srv.ServerInfo().Name() {
			return it, nil
		}
	}

	return nil, fmt.Errorf("item not found")
}
