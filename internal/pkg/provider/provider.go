package provider

import "go.minekube.com/gate/pkg/edition/java/proxy"

// Provider manipulates with registered servers.
type Provider interface {
	Name() string                        // Name returns providers name.
	Init() error                         // Init initializes the provider.
	Pause(proxy.RegisteredServer) error  // Pause pauses server.
	Resume(proxy.RegisteredServer) error // Resume resumes server.
}
