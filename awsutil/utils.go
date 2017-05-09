package awsutil

import (
	"fmt"

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
// and should pull it's configuration directly from the role on that
// instance.
//
// Alternatively, the following environment variables can be used to
// override the default session settings:
// AWS_ACCESS_KEY_ID
// AWS_SECRET_KEY_ID
// AWS_REGION

// StorkAWSClient contains the initialized clients needed for Stork
// to interact with AWS.
type StorkAWSClient struct {
	Config  *aws.Config
	Session *session.Session
	S3      s3iface.S3API
	EC2     ec2iface.EC2API
}

// NewStorkAWSClient returns a pointer to an initialized StorkAWSClient
func NewStorkAWSClient(config *config.StorkConfig) *StorkAWSClient {
	awsSession := session.Must(session.NewSession())
	logger.Debug("Created StorkAWSClient in region ", *awsSession.Config.Region)

	// Log every request made and its payload when debugging
	if config.Debug {
		awsSession.Handlers.Send.PushFront(func(r *request.Request) {
			logger.Debug(fmt.Sprintf(
				"AWS API: Request: %s/%s, Payload: %s",
				r.ClientInfo.ServiceName, r.Operation, r.Params,
			))
		})
	}

	return &StorkAWSClient{
		Session: awsSession,
		S3:      s3.New(awsSession),
		EC2:     ec2.New(awsSession),
	}
}

// CreateBucket creates a new S3 bucket with a given name
func (s *StorkAWSClient) CreateBucket(name string) error {
	logger.Debug("Creating bucket " + name)

	params := &s3.CreateBucketInput{
		Bucket: aws.String(name),
	}
	resp, err := s.S3.CreateBucket(params)

	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Debug("Created bucket at location: " + *resp.Location)
	return nil
}

// DeleteBucket deletes an existing S3 bucket and its contents, by name
func (s *StorkAWSClient) DeleteBucket(name string) error {
	logger.Warning("Deleting bucket " + name + " and its contents")

	params := &s3.DeleteBucketInput{
		Bucket: aws.String(name),
	}
	_, err := s.S3.DeleteBucket(params)

	if err != nil {
		logger.Error(err)
		return err
	}

	logger.Debug("Deleted bucket " + name)
	return nil
}
