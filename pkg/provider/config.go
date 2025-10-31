package provider

import (
	"time"

	"github.com/kasefuchs/lazygate/pkg/utils"
	"github.com/traefik/paerser/types"
)

// AllocationTimeConfig contains time related server configuration.
type AllocationTimeConfig struct {
	MinimumOnline       types.Duration // Minimum duration of time to stay online when server is starting.
	InactivityThreshold types.Duration // Duration of inactivity to stop after.
}

type AllocationConfig struct {
	Server    string                `validate:"required"` // The upstream server name.
	Namespace string                // Namespace to associate this allocation with.
	Queues    []string              // List of queues to try.
	Time      *AllocationTimeConfig // Time related server configuration.
}

func DefaultAllocationConfig() *AllocationConfig {
	return &AllocationConfig{
		Namespace: "default",
		Queues:    []string{"wait", "kick"},
		Time: &AllocationTimeConfig{
			MinimumOnline:       types.Duration(time.Minute),
			InactivityThreshold: types.Duration(time.Minute),
		},
	}
}

func ParseAllocationConfig(alloc Allocation) (*AllocationConfig, error) {
	cfg, err := alloc.ParseConfig(DefaultAllocationConfig(), utils.ChildLabel("allocation"))
	if err != nil {
		return nil, err
	}

	return cfg.(*AllocationConfig), nil
}
