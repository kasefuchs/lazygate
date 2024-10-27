package config

// Server represents upstream server config.
type Server struct {
	Name string // The name of the upstream server with which the service is associated.
}

type Config struct {
	Server Server // The upstream server config.
}
