package server

// DefaultConfig is the default set of configuration options for Stork.
var DefaultConfig = StorkConfig{
	ServerPort:   ":8080",
	DatabaseHost: "localhost:27017",
	DatabaseName: "stork",
	Debug:        false,
}

// StorkConfig encapsulates all Stork configuration options.
type StorkConfig struct {
	ServerPort   string
	DatabaseHost string
	DatabaseName string
	Debug        bool
}
