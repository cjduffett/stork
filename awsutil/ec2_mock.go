package awsutil

import "github.com/aws/aws-sdk-go/service/ec2/ec2iface"
import "github.com/aws/aws-sdk-go/service/ec2"

// EC2Mock mocks out the AWS EC2 API for testing
type EC2Mock struct {
	ec2iface.EC2API
	instances instanceMap
}

type instanceMock struct {
	state string
}

type instanceMap map[string]instanceMock

// NewEC2Mock returns a pointer to an initialized EC2 mock
func NewEC2Mock() *EC2Mock {
	return &EC2Mock{instances: make(instanceMap)}
}

// RunInstances mocks the ec2.runInstances operation
func (e *EC2Mock) RunInstances(*ec2.RunInstancesInput) (*ec2.Reservation, error) {
	// expected input includes:
	// ImageId, InstanceType, MinCount, MaxCount, SecurityGroupIds, SubnetId,
	// IamInstanceProfile, UserData (optional)
	return &ec2.Reservation{}, nil
}

// CreateTags mocks the ec2.createTags operation
func (e *EC2Mock) CreateTags(*ec2.CreateTagsInput) (*ec2.CreateTagsOutput, error) {
	return &ec2.CreateTagsOutput{}, nil
}

// TerminateInstances mocks the ec2.terminateInstances operation
func (e *EC2Mock) TerminateInstances(*ec2.TerminateInstancesInput) (*ec2.TerminateInstancesOutput, error) {
	return &ec2.TerminateInstancesOutput{}, nil
}

// DescribeInstanceStatus mocks the ec2.describeInstanceStatus operation
func (e *EC2Mock) DescribeInstanceStatus(*ec2.DescribeInstanceStatusInput) (*ec2.DescribeInstanceStatusOutput, error) {
	return &ec2.DescribeInstanceStatusOutput{}, nil
}
