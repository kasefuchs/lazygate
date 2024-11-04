package allocation

import (
	"time"

	ctypes "github.com/kasefuchs/lazygate/pkg/config/types"
	ptypes "github.com/traefik/paerser/types"
)

// Config represents specific server config.
type Config struct {
	Server string `validate:"required"` // The upstream server name.

	Time              *Time              // Time related server configuration.
	DisconnectReasons *DisconnectReasons // Reasons to kick with when disconnecting players.
}

// Time contains time related server configuration.
type Time struct {
	MinimumOnline       ptypes.Duration // Minimum duration of time to stay online when server is starting.
	InactivityThreshold ptypes.Duration // Duration of inactivity to stop after.
}

// DisconnectReasons contains reasons to use when disconnecting players.
type DisconnectReasons struct {
	Starting    ctypes.RawTextComponent // Reason to disconnect with when allocation is starting.
	StartFailed ctypes.RawTextComponent // Reason to disconnect with when action on allocation is failed.
}

// DefaultConfig returns default config.
func DefaultConfig() *Config {
	return &Config{
		Time: &Time{
			MinimumOnline:       ptypes.Duration(time.Minute),
			InactivityThreshold: ptypes.Duration(time.Minute),
		},
		DisconnectReasons: &DisconnectReasons{
			Starting:    "Server is starting...\n\nPlease try to reconnect in a minute.",
			StartFailed: "Failed to start the server.\n\nPlease try to reconnect later.",
		},
	}
}
