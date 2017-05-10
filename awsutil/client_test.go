package awsutil

import (
	"testing"

	"github.com/cjduffett/stork/config"
	"github.com/cjduffett/stork/logger"
	"github.com/stretchr/testify/suite"
)

type AWSUtilsTestSuite struct {
	suite.Suite
}

func TestAWSUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(AWSUtilsTestSuite))
}

func (a *AWSUtilsTestSuite) SetupSuite() {
	// verbose logging
	logger.LogLevel = logger.DebugLevel
}

func (a *AWSUtilsTestSuite) TestCreateBucket() {
	var err error
	client := newMockAWSClient()

	// Make a valid request
	err = client.CreateBucket("test-bucket")
	a.NoError(err)

	// Creating a bucket that already exists should fail
	err = client.CreateBucket("test-bucket")
	a.Error(err)
}

func (a *AWSUtilsTestSuite) TestDeleteBucket() {
	var err error
	client := newMockAWSClient()

	// Create a bucket
	err = client.CreateBucket("test-bucket")
	a.NoError(err)

	// Make a valid request to delete that bucket
	err = client.DeleteBucket("test-bucket")
	a.NoError(err)

	// Trying to delete a bucket that doesn't exist should fail
	err = client.DeleteBucket("foo-bucket")
	a.Error(err)
}

func newMockAWSClient() *AWSClient {
	return &AWSClient{
		Config:  config.DefaultConfig,
		Session: nil,
		S3:      NewS3Mock(),
		EC2:     NewEC2Mock(),
	}
}
