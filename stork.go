package main

import (
	"flag"

	"github.com/cjduffett/stork/server"
)

func main() {
	debug := flag.Bool("debug", server.DefaultConfig.Debug, "Enable debug level logging")
	flag.Parse()

	config := server.DefaultConfig
	config.Debug = *debug

	s := server.NewServer(config)
	s.Run()
}
