package dynamic

// Config represents specific server config.
type Config struct {
	Server string `validate:"required"` // The upstream server name.
}
