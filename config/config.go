package config

// DefaultConfig is the default set of configuration options for Stork.
var DefaultConfig = &StorkConfig{
	ServerHost:     "",
	ServerPort:     "8080",
	StaticFilePath: "assets",
	DatabaseHost:   "localhost:27017",
	DatabaseName:   "stork",
	Debug:          false,
	DoneEndpoint:   "/task/:id/done",
}

// StorkConfig encapsulates all Stork configuration options.
type StorkConfig struct {
	// Server configuration
	ServerHost     string
	ServerPort     string
	StaticFilePath string
	Debug          bool

	// Database configuration
	DatabaseHost string
	DatabaseName string

	// AWS Configuration
	DoneEndpoint   string
	SyntheaImageID string
}
