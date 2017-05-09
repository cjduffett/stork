package awsutil

import (
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

	// If it isn't, or is only partially set, warn the user
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

// DescribeBucket returns the status of an S3 bucket, by name.
func (s *AWSClient) DescribeBucket(name string) (*BucketStatus, error) {
	return nil, nil
}

// StartInstances starts n new Synthea instances with the same configuration.
func (s *AWSClient) StartInstances(n int, config InstanceConfig) ([]string, error) {
	// Make a RunInstances request for n instances
	// returns []Instance
	// Don't wait for the instances to start/confirm starting

	// Tag these instances

	// Return []InstanceIDs
	return nil, nil
}

// TerminateInstances terminates one or more Synthea instances.
// This may be called after the /done endpoint is pinged, or if
// an abort request is made.
func (s *AWSClient) TerminateInstances(instanceIDs []string) error {
	return nil
}

// DescribeInstances returns the status of one or more Synthea instances.
func (s *AWSClient) DescribeInstances(instanceIDs []string) ([]InstanceStatus, error) {
	return nil, nil
}
