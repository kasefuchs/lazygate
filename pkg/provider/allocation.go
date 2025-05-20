package provider

// AllocationState represents allocation state.
type AllocationState uint8

const (
	AllocationStateUnknown AllocationState = iota // Allocation is in unknown state.
	AllocationStateStarted                        // Allocation is in started state.
	AllocationStateStopped                        // Allocation is in stopped state.
)

// Allocation represents physical allocation.
type Allocation interface {
	Stop() error  // Stop stops the allocation.
	Start() error // Start starts the allocation.

	State() AllocationState                                          // State returns current allocation state.
	ParseConfig(cfg interface{}, prefix string) (interface{}, error) // ParseConfig parses specific config of allocation.
}
