package provider

import "go.minekube.com/gate/pkg/edition/java/proxy"

// InitOptions represents options to pass to provider.
type InitOptions struct{}

// Provider manipulates with registered servers.
type Provider interface {
	Init(opts InitOptions) error // Init initializes the provider.

	AllocationGet(srv proxy.RegisteredServer) (Allocation, error) // GetAllocation returns matching Allocation.
	AllocationList() ([]Allocation, error)                        // GetAll returns all matching Allocation's.
}
