package server

import (
	"net/http/httptest"
	"testing"

	"github.com/cjduffett/stork/api"
	"github.com/cjduffett/stork/config"
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
	storkServer.Session = s.DB().Session

	// Just testing the middleware in this package. For API and site tests, see the
	// "api" and "site" packages, respectively.
	api.RegisterRoutes(storkServer.Engine, storkServer.Session, config)
	s.StorkServer = httptest.NewServer(storkServer.Engine)
}

func (s *ServerTestSuite) TearDownSuite() {
	s.StorkServer.Close()
	// Clean up and remove all temporary files from the mocked database.
	// See testutil/mongo_suite.go for more.
	s.TearDownDBServer()
}
