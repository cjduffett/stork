package server

import (
	"net/http/httptest"
	"testing"

	"github.com/cjduffett/stork/api"
	"github.com/cjduffett/stork/awsutil"
	"github.com/cjduffett/stork/config"
	"github.com/cjduffett/stork/db"
	"github.com/cjduffett/stork/testutil"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	testutil.MongoSuite
	StorkServer *httptest.Server
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (s *ServerTestSuite) SetupSuite() {
	// Set gin to release mode (less verbose output)
	gin.SetMode(gin.ReleaseMode)

	// Create a mock StorkServer
	config := config.DefaultConfig
	config.DatabaseName = "stork-test"
	storkServer := NewServer(config)

	// Create a mock MongoDB connection
	storkServer.Session = s.DB().Session

	// Register middleware (CORS, etc.)
	RegisterMiddleware(storkServer.Engine)

	// Create a new Data Access Layer
	dal := db.NewDataAccessLayer(storkServer.Session, config.DatabaseName)

	// Create a new AWSClient
	awsClient := awsutil.NewAWSClient(config.Debug)

	// Register API routes and setup controllers
	api.RegisterRoutes(storkServer.Engine, dal, awsClient)

	// Start the httptest server
	s.StorkServer = httptest.NewServer(storkServer.Engine)
}

func (s *ServerTestSuite) TearDownSuite() {
	s.StorkServer.Close()
	// Clean up and remove all temporary files from the mocked database.
	// See testutil/mongo_suite.go for more.
	s.TearDownDBServer()
}
