package allocation

import (
	"time"

	ctypes "github.com/kasefuchs/lazygate/pkg/config/types"
	ptypes "github.com/traefik/paerser/types"
)

// Config represents specific server config.
type Config struct {
	Server string `validate:"required"` // The upstream server name.

	Time  *Time  // Time related server configuration.
	Queue *Queue // Queue related configuration.
}

// Time contains time related server configuration.
type Time struct {
	MinimumOnline       ptypes.Duration // Minimum duration of time to stay online when server is starting.
	InactivityThreshold ptypes.Duration // Duration of inactivity to stop after.
}

// Queue configuration.
type Queue struct {
	Try  []string   // List of queues to try.
	Wait *QueueWait // Queue waiting for server start.
	Kick *QueueKick // Queue instantly kicking player.
}

// QueueKick represents kick queue configuration.
type QueueKick struct {
	Starting ctypes.RawTextComponent // Reason to kick with when allocation is starting.
}

// QueueWait represents wait queue configuration.
type QueueWait struct {
	Timeout      ptypes.Duration // Maximum amount of time to wait for server start. Must be less than Minecraft timeout of 30 seconds.
	PingInterval ptypes.Duration // Interval to ping backend server with. Keep it lower than half of timeout.
}

// DefaultConfig returns default config.
func DefaultConfig() *Config {
	return &Config{
		Time: &Time{
			MinimumOnline:       ptypes.Duration(time.Minute),
			InactivityThreshold: ptypes.Duration(time.Minute),
		},
		Queue: &Queue{
			Try: []string{"wait", "kick"},
			Kick: &QueueKick{
				Starting: "Server is starting...\n\nPlease try to reconnect in a minute.",
			},
			Wait: &QueueWait{
				Timeout:      ptypes.Duration(25 * time.Second),
				PingInterval: ptypes.Duration(3 * time.Second),
			},
		},
	}
}
