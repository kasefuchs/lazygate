package plugin

import (
	"fmt"
	"os"
	"strings"

	"github.com/kasefuchs/lazygate/pkg/provider"
	"github.com/kasefuchs/lazygate/pkg/provider/docker"
	"github.com/kasefuchs/lazygate/pkg/provider/nomad"
)

// providerSelector
type providerSelector func() (provider.Provider, error)

// Options contains customizable plugin options.
type Options struct {
	ProviderSelector providerSelector // Selector of provider.
}

// DefaultProviderSelector contains default provider selector.
func DefaultProviderSelector() (provider.Provider, error) {
	name := os.Getenv("LAZYGATE_PROVIDER")

	switch strings.TrimSpace(name) {
	case "nomad":
		return &nomad.Provider{}, nil
	case "docker":
		return &docker.Provider{}, nil
	case "":
		return nil, fmt.Errorf("no plugin provider specified")
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}

// DefaultOptions returns options object with default parameters.
func DefaultOptions() *Options {
	return &Options{
		ProviderSelector: DefaultProviderSelector,
	}
}
