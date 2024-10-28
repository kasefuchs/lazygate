package provider

import (
	"fmt"
	"strings"

	"github.com/kasefuchs/lazygate/internal/pkg/provider/docker"
	"github.com/kasefuchs/lazygate/internal/pkg/provider/nomad"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

// Provider manipulates with registered servers.
type Provider interface {
	Init() error                          // Init initializes the provider.
	Pause(*proxy.RegisteredServer) error  // Pause pauses server.
	Resume(*proxy.RegisteredServer) error // Resume resumes server.
}

// NewProvider returns new Provider by its name.
func NewProvider(name string) (Provider, error) {
	switch strings.ToLower(name) {
	case "nomad":
		return &nomad.Provider{}, nil
	case "docker":
		return &docker.Provider{}, nil
	case "":
		return nil, fmt.Errorf("no provider specified")
	default:
		return nil, fmt.Errorf("unknown provider %s", name)
	}
}
