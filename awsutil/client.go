package awsutil

import (
	b64 "encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/cjduffett/stork/config"
	"github.com/cjduffett/stork/db"
	"github.com/cjduffett/stork/logger"
)

// This file implements a series of utilities that simplify working
// with AWS resources. Stork will ultimately be run on an EC2 instance
// and should pull it's configuration directly from the environment
// on that instance.

// The following environment variables must be set:
const (
	requiredAccessKeyEnvVar = "AWS_ACCESS_KEY_ID"
	requiredSecretKeyEnvVar = "AWS_SECRET_ACCESS_KEY"
	requiredRegionEnvVar    = "AWS_REGION"
)

// AWSClient contains the initialized clients and interfaces
// needed for Stork to interact with AWS.
type AWSClient struct {
	Config  *config.StorkConfig
	Session *session.Session
	S3      s3iface.S3API
	EC2     ec2iface.EC2API
}

// NewAWSClient returns a pointer to an initialized AWSClient
func NewAWSClient(config *config.StorkConfig) *AWSClient {

	// Check that the environment is set
	required := []string{requiredAccessKeyEnvVar, requiredSecretKeyEnvVar, requiredRegionEnvVar}
	missing := []string{}
	for _, r := range required {
		if os.Getenv(r) == "" {
			missing = append(missing, r)
		}
	}

	// If it isn't, or is only partially set, warn the user. If Stork is running
	// on an EC2 instance it may be able to pull its configuration from the machine.
	if len(missing) > 0 {
		logger.Warning("One or more environment variables are not set, this may prevent Stork from connecting to AWS")
		for _, m := range missing {
			logger.Warning("Missing environment variable " + m)
		}
	}

	// Establish a session with AWS
	awsSession := session.Must(session.NewSession())
	region := *awsSession.Config.Region
	if region == "" {
		region = "UNKNOWN"
	}
	logger.Info("Connecting Stork to AWS in region " + region)

	// When debugging, log every request made and its payload
	if config.Debug {
		awsSession.Handlers.Send.PushFront(func(r *request.Request) {
			logger.Debug(fmt.Sprintf(
				"AWS API: Request: %s/%s, Payload: %s",
				r.ClientInfo.ServiceName, r.Operation, r.Params,
			))
		})
	}

	return &AWSClient{
		Config:  config,
		Session: awsSession,
		S3:      s3.New(awsSession),
		EC2:     ec2.New(awsSession),
	}
}

// CreateBucket creates a new S3 bucket with a given name
func (s *AWSClient) CreateBucket(name string) error {
	logger.Debug("Creating bucket " + name)

	params := &s3.CreateBucketInput{
		Bucket: aws.String(name),
	}
	resp, err := s.S3.CreateBucket(params)

	if err != nil {
		logger.Error("Failed to create bucket " + name)
		return err
	}

	logger.Debug("Created bucket at location: " + *resp.Location)
	return nil
}

// DeleteBucket deletes an existing S3 bucket and its contents, by name
func (s *AWSClient) DeleteBucket(name string) error {
	logger.Debug("Deleting bucket " + name + " and its contents")

	params := &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	}
	_, err := s.S3.DeleteBucket(params)

	if err != nil {
		logger.Error("Failed to delete bucket " + name)
		return err
	}

	logger.Debug("Deleted bucket " + name)
	return nil
}

// StartInstances starts n new Synthea instances with the same configuration.
// All instances are expected to share an equal compute load, with a minimum
// of 500 patients each (this is validated elsewhere).
func (s *AWSClient) StartInstances(n int64, iConfig *InstanceConfig) ([]string, error) {
	var err error

	logger.Debug(fmt.Sprintf("Starting %d instances of Synthea for task %s", n, iConfig.TaskID))

	// The InstanceConfig must be validated before doing anything.
	if !ValidateConfig(iConfig, s.Config) {
		return nil, errors.New("Invalid InstanceConfig")
	}

	// Serialie the InstanceConfig to pass it as UserData.
	rawUserData, err := json.Marshal(iConfig)
	if err != nil {
		return nil, err
	}

	// The raw user data must be base64 encoded. It will be automatically
	// decoded when Synthea requests it from the EC2 instance user data.
	encodedUserData := b64.StdEncoding.EncodeToString(rawUserData)

	// Make a RunInstances request for n Synthea instances
	runParams := &ec2.RunInstancesInput{
		ImageId:          aws.String(s.Config.SyntheaImageID),
		InstanceType:     aws.String(s.Config.SyntheaInstanceType),
		MinCount:         aws.Int64(n),
		MaxCount:         aws.Int64(n),
		SecurityGroupIds: []*string{aws.String(s.Config.SyntheaSecurityGroupID)},
		IamInstanceProfile: &ec2.IamInstanceProfileSpecification{
			// This is an ARN to an EC2 Role with one or more associated policies
			Arn: aws.String(s.Config.SyntheaRoleArn),
		},
		SubnetId: aws.String(s.Config.SyntheaSubnetID),
		UserData: aws.String(encodedUserData),
	}
	reservation, err := s.EC2.RunInstances(runParams)
	if err != nil {
		logger.Error("Failed to start instances for task " + iConfig.TaskID)
		return nil, err
	}

	// Parse the instance IDs out of the reservation
	instanceIDs := make([]*string, len(reservation.Instances))
	for i, instance := range reservation.Instances {
		instanceIDs[i] = instance.InstanceId
	}
	strInstanceIDs := toStrings(instanceIDs)

	logger.Debug(fmt.Sprintf("Started %d instances: %v", n, strInstanceIDs))

	// Tag these instance with "stork-synthea" and the taskID
	// so we know who they belong to
	logger.Debug(fmt.Sprintf("Tagging instances %v", strInstanceIDs))

	tagParams := &ec2.CreateTagsInput{
		Resources: instanceIDs,
		Tags: []*ec2.Tag{
			&ec2.Tag{
				Key:   aws.String("role"),
				Value: aws.String("stork-synthea"),
			},
			&ec2.Tag{
				Key:   aws.String("task"),
				Value: aws.String(iConfig.TaskID),
			},
		},
	}
	_, err = s.EC2.CreateTags(tagParams)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to tag instances %v", strInstanceIDs))
		return nil, err
	}
	logger.Debug(fmt.Sprintf("Tagged instances %v", strInstanceIDs))

	return strInstanceIDs, nil
}

// TerminateInstances terminates one or more Synthea instances.
// This may be called after the /done endpoint is pinged, or if
// an abort request is made.
func (s *AWSClient) TerminateInstances(instanceIDs []string) error {
	logger.Debug(fmt.Sprintf("Terminating instances %v", instanceIDs))
	params := &ec2.TerminateInstancesInput{
		InstanceIds: toAWSStrings(instanceIDs),
	}
	_, err := s.EC2.TerminateInstances(params)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to terminate instances %v", instanceIDs))
		return err
	}
	return nil
}

// DescribeInstanceStatus returns the status of one or more Synthea instances.
// The status returned from ec2.DescribeInstanceStatus is converted to a local
// representation of status.
func (s *AWSClient) DescribeInstanceStatus(instanceIDs []string) ([]InstanceStatus, error) {
	logger.Debug(fmt.Sprintf("Getting status of instances %v", instanceIDs))
	params := &ec2.DescribeInstanceStatusInput{
		InstanceIds: toAWSStrings(instanceIDs),
	}
	resp, err := s.EC2.DescribeInstanceStatus(params)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to get status of instances %v", instanceIDs))
		return nil, err
	}

	// Parse the response into our own internal representation of instance status
	statuses := make([]InstanceStatus, len(instanceIDs))
	for i, status := range resp.InstanceStatuses {
		statuses[i] = InstanceStatus{
			InstanceID: *status.InstanceId,
			Status:     convertInstanceStatus(status.InstanceState),
		}
	}
	return statuses, nil
}

// Converts and AWS instance state to a locally known string value.
// For our purposes, any instance that isn't terminated is considered active.
func convertInstanceStatus(state *ec2.InstanceState) string {
	status := ""
	switch *state.Name {
	case ec2.InstanceStateNameTerminated:
		status = db.InstanceStatusDone
	default:
		status = db.InstanceStatusActive
	}
	return status
}

// Converts an array of AWS strings (*string) to ordinary strings
func toStrings(ptrs []*string) []string {
	out := make([]string, len(ptrs))
	for i, ptr := range ptrs {
		out[i] = *ptr
	}
	return out
}

// Converts an array of ordinary strings to AWS strings (*string)
func toAWSStrings(strs []string) []*string {
	out := make([]*string, len(strs))
	for i, str := range strs {
		out[i] = aws.String(str)
	}
	return out
}
