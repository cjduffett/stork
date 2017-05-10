package main

import (
	"flag"

	"github.com/cjduffett/stork/config"
	"github.com/cjduffett/stork/server"
)

func main() {
	conf := config.DefaultConfig

	// Server options
	host := flag.String("host", config.DefaultConfig.ServerHost, "StorkServer host")
	port := flag.String("port", config.DefaultConfig.ServerPort, "StorkServer port")
	debug := flag.Bool("debug", config.DefaultConfig.Debug, "Enable debug level logging")

	// Database options - all database options begin with "db."
	dbhost := flag.String("db.host", config.DefaultConfig.DatabaseHost, "Database host")
	dbname := flag.String("db.name", config.DefaultConfig.DatabaseName, "Database name")

	// AWS options - all aws options begin with "aws."
	syntheaImageID := flag.String("aws.synthea-image-id", config.DefaultConfig.SyntheaImageID, "The Synthea AMI ID to run")
	syntheaInstanceType := flag.String("aws.synthea-instance-type", config.DefaultConfig.SyntheaInstanceType, "The type of EC2 instance to run Synthea on")
	syntheaSecurityGroupID := flag.String("aws.synthea-security-group-id", config.DefaultConfig.SyntheaSecurityGroupID, "The security group associated with a Synthea instance")
	syntheaRoleArn := flag.String("aws.synthea-role-arn", config.DefaultConfig.SyntheaRoleArn, "The role associated with a Synthea instance")

	flag.Parse()
	conf.ServerHost = *host
	conf.ServerPort = *port
	conf.Debug = *debug

	conf.DatabaseHost = *dbhost
	conf.DatabaseName = *dbname

	conf.SyntheaImageID = *syntheaImageID
	conf.SyntheaInstanceType = *syntheaInstanceType
	conf.SyntheaSecurityGroupID = *syntheaSecurityGroupID
	conf.SyntheaRoleArn = *syntheaRoleArn

	s := server.NewServer(conf)
	s.Run()
}
