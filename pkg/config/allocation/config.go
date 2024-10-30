package allocation

// Config represents specific server config.
type Config struct {
	Server            string             `validate:"required"` // The upstream server name.
	DisconnectReasons *DisconnectReasons // Reasons to kick with when disconnecting players.
}

// DisconnectReasons contains reasons to use when disconnecting players.
type DisconnectReasons struct {
	Starting     string // Reason to disconnect with when allocation is starting.
	ActionFailed string // Reason to disconnect with when action on allocation is failed.
}

// DefaultConfig returns default config.
func DefaultConfig() *Config {
	return &Config{
		DisconnectReasons: &DisconnectReasons{
			Starting:     "Server is starting...\n\nPlease try to reconnect in a minute.",
			ActionFailed: "The action on the allocation failed.\n\nPlease try to reconnect later.",
		},
	}
}
