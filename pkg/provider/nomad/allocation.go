package nomad

import (
	"fmt"

	"github.com/hashicorp/nomad/api"
	"github.com/kasefuchs/lazygate/pkg/provider"
)

var _ provider.Allocation = (*Allocation)(nil)

var (
	scaleStop  int64 = 0
	scaleStart int64 = 1
)

// Allocation represents Nomad provider item.
type Allocation struct {
	client *api.Client // Nomad API client.
	item   *item       // Nomad task item.
}

func NewAllocation(client *api.Client, item *item) *Allocation {
	return &Allocation{
		client: client,
		item:   item,
	}
}

func (a *Allocation) Stop() error {
	return a.scale(&scaleStop)
}

func (a *Allocation) Start() error {
	return a.scale(&scaleStart)
}

func (a *Allocation) scale(count *int64) error {
	job, _, err := a.client.Jobs().Info(*a.item.job.ID, nil)
	if err != nil {
		return err
	}

	group := job.LookupTaskGroup(*a.item.group.Name)
	if group == nil {
		return fmt.Errorf("task group %s not found", *a.item.group.Name)
	}

	if int64(*group.Count) != *count {
		req := &api.ScalingRequest{
			Count: count,
			Target: map[string]string{
				"Job":   *a.item.job.ID,
				"Group": *a.item.group.Name,
			},
			JobModifyIndex: 0,
			PolicyOverride: false,
		}
		if _, _, err := a.client.Jobs().ScaleWithRequest(*a.item.job.ID, req, nil); err != nil {
			return err
		}
	}

	return nil
}
