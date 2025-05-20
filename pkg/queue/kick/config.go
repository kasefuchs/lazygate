package kick

// TicketConfig represents kick queue configuration.
type TicketConfig struct {
	Reason RawTextComponent // Reason to kick with when allocation is starting.
}
