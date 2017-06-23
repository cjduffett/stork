package config

// DefaultConfig is the default set of configuration options for Stork.
// Note: with this default configuration Stork has enough information to start,
// but not to make requests to AWS. Those configuration options will need
// to be passed from the command line.
var DefaultConfig = &StorkConfig{
	ServerHost: "localhost",
	ServerPort: "8080",
	Debug:      false,

	DatabaseHost: "localhost:27017",
	DatabaseName: "stork",

	SyntheaImageID:         "",
	SyntheaInstanceType:    "t2.micro",
	SyntheaSecurityGroupID: "",
	SyntheaRoleArn:         "",
	SyntheaSubnetID:        "subnet-581c7275",

	MinPopulationSize: 500,
	DoneEndpoint:      "/task/:id/done",
}

// StorkConfig encapsulates all Stork configuration options.
type StorkConfig struct {
	// Stork Server configuration options.
	ServerHost string
	ServerPort string
	Debug      bool

	// MongoDB configuration options.
	DatabaseHost string
	DatabaseName string

	// The prebuilt Snythea image (already available in AWS) to use.
	SyntheaImageID string

	// The EC2 instance type to run Synthea on. A compute-optimized instance
	// is recommended, for example a c4.large instance.
	SyntheaInstanceType string

	// The security groups ID (already created in AWS) for Snythea to use.
	// Minimally, Synthea must be able to make outbound HTTP/HTTPS requests.
	SyntheaSecurityGroupID string

	// The ARN of the IAM role (already created in AWS) that Synthea should use.
	// Minimally, Synthea must be able to perform an s3:putObject action.
	SyntheaRoleArn string

	// The VPC subnet to run Synthea in. Typically this a private subnet that must
	// be in the same region as Stork and the S3 bucket Synthea writes to.
	SyntheaSubnetID string

	// The minimum number of patient records that an instance should generate.
	// This is practically defined by the towns.json data the feeds a sequential
	// synthea run. In that case, 481 patients are generated.
	MinPopulationSize int

	// The stork endpoint that Synthea instances should ping when done.
	DoneEndpoint string
}
