package main

import (
	"flag"

	"github.com/cjduffett/stork/config"
	"github.com/cjduffett/stork/server"
)

func main() {
	conf := config.DefaultConfig

	// Server options
	port := flag.String("port", config.DefaultConfig.ServerPort, "StorkServer port")
	debug := flag.Bool("debug", config.DefaultConfig.Debug, "Enable debug level logging")

	// Database options - all database options begin with "db"
	dbhost := flag.String("db.host", config.DefaultConfig.DatabaseHost, "Database host")
	dbname := flag.String("db.name", config.DefaultConfig.DatabaseName, "Database name")

	flag.Parse()
	conf.ServerPort = *port
	conf.Debug = *debug
	conf.DatabaseHost = *dbhost
	conf.DatabaseName = *dbname

	s := server.NewServer(conf)
	s.Run()
}
