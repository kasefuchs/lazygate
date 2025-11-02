package send

import "github.com/traefik/paerser/types"

// TicketConfig represents send queue configuration.
type TicketConfig struct {
	PingInterval types.Duration // Interval to ping backend server with.
	Timeout      types.Duration // Maximum amount of time to wait for server start.
	To           string         // Temporary server to send the player while the target starts.
}
