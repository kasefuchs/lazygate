package nomad

import (
	"github.com/hashicorp/nomad/api"
	"github.com/kasefuchs/lazygate/pkg/config/allocation"
)

// Allocation internal data in Nomad context.
type item struct {
	job   *api.Job       // Task job
	group *api.TaskGroup // Task group
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
				if _, err := allocation.ParseTags(service.Tags); err != nil {
					continue
				}

				items = append(items, &item{
					job:   job,
					group: taskGroup,
				})
			}
		}
	}

	return items, nil
}
