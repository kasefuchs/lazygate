package provider

// Allocation represents physical allocation.
type Allocation interface {
	Stop() error  // Stop stops the item.
	Start() error // Start starts the item.
}
