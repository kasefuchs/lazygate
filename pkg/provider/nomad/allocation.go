package nomad

import (
	"fmt"

	"github.com/hashicorp/nomad/api"
	"github.com/kasefuchs/lazygate/pkg/provider"
	"github.com/kasefuchs/lazygate/pkg/utils"
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

func (a *Allocation) taskGroup() (*api.TaskGroup, error) {
	job, _, err := a.client.Jobs().Info(*a.item.job.ID, nil)
	if err != nil {
		return nil, err
	}

	group := job.LookupTaskGroup(*a.item.group.Name)
	if group == nil {
		return nil, fmt.Errorf("task group %s not found", *a.item.group.Name)
	}

	return group, nil
}

func (a *Allocation) scale(count *int64) error {
	group, err := a.taskGroup()
	if err != nil {
		return err
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

func (a *Allocation) Stop() error {
	return a.scale(&scaleStop)
}

func (a *Allocation) Start() error {
	return a.scale(&scaleStart)
}

func (a *Allocation) State() provider.AllocationState {
	group, err := a.taskGroup()
	if err != nil {
		return provider.AllocationStateUnknown
	}

	if *group.Count == 0 {
		return provider.AllocationStateStopped
	}

	return provider.AllocationStateStarted
}

func (a *Allocation) ParseConfig(cfg interface{}, rootLabel string) (interface{}, error) {
	group, err := a.taskGroup()
	if err != nil {
		return nil, err
	}

	for _, service := range group.Services {
		cfg, err := utils.ParseTags(service.Tags, cfg, rootLabel)
		if err != nil {
			continue
		}

		return cfg.(*provider.AllocationConfig), nil
	}

	return nil, fmt.Errorf("no allocation configuration found")
}
