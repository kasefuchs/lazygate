package provider

import (
	"context"

	"go.minekube.com/gate/pkg/edition/java/proxy"
)

const LogName = "lazygate.provider"

// InitOptions represents options to pass to provider initializer.
type InitOptions struct {
	Ctx    context.Context // Plugin context.
	Config interface{}     // Provider config.
}

// Provider manipulates with registered servers.
type Provider interface {
	Init(opts *InitOptions) error // Init initializes the provider.
	Name() string                 // Name returns name of provider.
	DefaultConfig() interface{}   // DefaultConfig returns default plugin config.

	AllocationGet(srv proxy.RegisteredServer) (Allocation, error) // AllocationGet returns matching Allocation.
	AllocationList() ([]Allocation, error)                        // AllocationList returns all matching Allocation's.
}
