package wait

import "github.com/traefik/paerser/types"

// TicketConfig represents wait queue configuration.
type TicketConfig struct {
	PingInterval types.Duration // Interval to ping backend server with. Keep it lower than half of timeout.
	Timeout      types.Duration // Maximum amount of time to wait for server start. Must be less than Minecraft timeout of 30 seconds.
}
