package awsutil

import (
	"testing"

	"github.com/cjduffett/stork/config"
	"github.com/stretchr/testify/suite"
)

type TypesTestSuite struct {
	suite.Suite
	AWSClient *AWSClient
}

func TestTypesTestSuite(t *testing.T) {
	suite.Run(t, new(TypesTestSuite))
}

func (t *TypesTestSuite) TestValidateInstanceConfig() {
	// Validation also requires the Stork configuration object
	sConfig := config.DefaultConfig

	// First, make sure a valid config returns true
	iConfig := &InstanceConfig{
		TaskID:       "123abc",
		Population:   sConfig.MinPopulationSize + 100,
		BucketName:   "123abc-bucket",
		DoneEndpoint: "https://stork.com/tasks/:id/done",
	}
	t.True(ValidateConfig(iConfig, sConfig))

	// Now test with an invalid population number (less than config.MinPopulationSize)
	iConfig.Population = sConfig.MinPopulationSize - 2
	t.False(ValidateConfig(iConfig, sConfig))

	// Now test with an empty field
	iConfig.Population = sConfig.MinPopulationSize + 1000
	iConfig.TaskID = ""
	t.False(ValidateConfig(iConfig, sConfig))
}
