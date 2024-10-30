package provider

import "github.com/kasefuchs/lazygate/pkg/config/allocation"

// Allocation represents physical allocation.
type Allocation interface {
	Stop() error  // Stop stops the allocation.
	Start() error // Start starts the allocation.

	Config() *allocation.Config // Config of allocation.
}
