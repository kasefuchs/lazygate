package queue

import (
	"github.com/kasefuchs/lazygate/pkg/registry"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

// InitOptions represents options to pass to queue initializer.
type InitOptions struct {
	Proxy *proxy.Proxy
}

// Ticket represents metadata to pass to queue.
type Ticket struct {
	Entry  *registry.Entry
	Config interface{}
	Player proxy.Player
}

// Queue manipulates players waiting for server allocation start.
type Queue interface {
	Name() string                     // Name of queue.
	Init(opts *InitOptions) error     // Initializes the queue.
	DefaultTicketConfig() interface{} // Returns default ticket config.

	Enter(ticket *Ticket) bool // Enters player to queue using ticket. Returns false if player required to enter next queue.
}
