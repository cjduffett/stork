package awsutil

import "github.com/aws/aws-sdk-go/service/ec2/ec2iface"

// EC2Mock mocks out the AWS EC2 API for testing
type EC2Mock struct {
	ec2iface.EC2API
}
