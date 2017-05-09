package awsutil

import (
	"testing"

	"github.com/cjduffett/stork/logger"
	"github.com/stretchr/testify/suite"
)

type AWSUtilsTestSuite struct {
	suite.Suite
	StorkAWSClient *StorkAWSClient
}

func TestAWSUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(AWSUtilsTestSuite))
}

func (a *AWSUtilsTestSuite) SetupSuite() {
	// verbose logging
	logger.LogLevel = logger.DebugLevel

	a.StorkAWSClient = &StorkAWSClient{
		Session: nil,
		S3:      NewS3Mock(),
		EC2:     EC2Mock{},
	}
}

func (a *AWSUtilsTestSuite) TestCreateBucket() {
	var err error

	// Make a valid request
	err = a.StorkAWSClient.CreateBucket("test-bucket")
	a.NoError(err)

	// Creating a bucket that already exists should fail
	err = a.StorkAWSClient.CreateBucket("test-bucket")
	a.Error(err)
}

func (a *AWSUtilsTestSuite) TestDeleteBucket() {

}
