package plugin

import (
	"github.com/kasefuchs/lazygate/pkg/provider"
	"github.com/kasefuchs/lazygate/pkg/provider/docker"
	"github.com/kasefuchs/lazygate/pkg/provider/nomad"
	"github.com/kasefuchs/lazygate/pkg/provider/pufferpanel"
	"github.com/kasefuchs/lazygate/pkg/queue"
	"github.com/kasefuchs/lazygate/pkg/queue/kick"
	"github.com/kasefuchs/lazygate/pkg/queue/send"
	"github.com/kasefuchs/lazygate/pkg/queue/wait"
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
