package plugin

import (
	"fmt"
	"os"
	"strings"

	pconfig "github.com/kasefuchs/lazygate/pkg/config/plugin"
	"github.com/kasefuchs/lazygate/pkg/provider"
	"github.com/kasefuchs/lazygate/pkg/provider/docker"
	"github.com/kasefuchs/lazygate/pkg/provider/nomad"
	"github.com/kasefuchs/lazygate/pkg/provider/pufferpanel"
	"github.com/kasefuchs/lazygate/pkg/queue"
	"github.com/kasefuchs/lazygate/pkg/queue/kick"
	"github.com/kasefuchs/lazygate/pkg/queue/wait"
)

// providerSelector represents function used to select provider to use.
type providerSelector func() (provider.Provider, error)

// queuesSelector represents function used to select queues to use.
type queuesSelector func() ([]queue.Queue, error)

// configLoader represents function used to load plugin configuration.
type configLoader func() (*pconfig.Config, error)

// Options contains customizable plugin options.
type Options struct {
	ProviderSelector providerSelector // Selector of provider.
	QueuesSelector   queuesSelector   // Selector of available queues.
	ConfigLoader     configLoader     // Loader of plugin config.
}

// DefaultProviderSelector contains default provider selector.
func DefaultProviderSelector() (provider.Provider, error) {
	name := os.Getenv("LAZYGATE_PROVIDER")

	switch strings.TrimSpace(name) {
	case "nomad":
		return &nomad.Provider{}, nil
	case "docker":
		return &docker.Provider{}, nil
	case "pufferpanel":
		return &pufferpanel.Provider{}, nil
	case "":
		return nil, fmt.Errorf("no allocation provider specified")
	default:
		return nil, fmt.Errorf("unknown allocation provider: %s", name)
	}
}

// DefaultQueuesSelector contains default queues selector.
func DefaultQueuesSelector() ([]queue.Queue, error) {
	return []queue.Queue{
		&wait.Queue{},
		&kick.Queue{},
	}, nil
}

// DefaultConfigLoader loads plugin config from environment.
func DefaultConfigLoader() (*pconfig.Config, error) {
	return pconfig.ParseEnv()
}

// DefaultOptions returns options object with default parameters.
func DefaultOptions() *Options {
	return &Options{
		ProviderSelector: DefaultProviderSelector,
		QueuesSelector:   DefaultQueuesSelector,
		ConfigLoader:     DefaultConfigLoader,
	}
}
