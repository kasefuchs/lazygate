package plugin

import (
	"codeberg.org/kasefuchs/lazygate/pkg/provider"
	"codeberg.org/kasefuchs/lazygate/pkg/provider/docker"
	"codeberg.org/kasefuchs/lazygate/pkg/provider/nomad"
	"codeberg.org/kasefuchs/lazygate/pkg/provider/pufferpanel"
	"codeberg.org/kasefuchs/lazygate/pkg/queue"
	"codeberg.org/kasefuchs/lazygate/pkg/queue/kick"
	"codeberg.org/kasefuchs/lazygate/pkg/queue/send"
	"codeberg.org/kasefuchs/lazygate/pkg/queue/wait"
)

// Options contains customizable plugin options.
type Options struct {
	Queues    []queue.Queue       // Registered queues.
	Providers []provider.Provider // Registered providers.
}

// DefaultOptions returns options object with default parameters.
func DefaultOptions() *Options {
	return &Options{
		Queues: []queue.Queue{
			&kick.Queue{},
			&send.Queue{},
			&wait.Queue{},
		},
		Providers: []provider.Provider{
			&docker.Provider{},
			&nomad.Provider{},
			&pufferpanel.Provider{},
		},
	}
}
