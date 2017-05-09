package server

import (
	"fmt"
	"os"
	"strings"

	"github.com/cjduffett/stork/api"
	"github.com/cjduffett/stork/awsutil"
	"github.com/cjduffett/stork/config"
	"github.com/cjduffett/stork/db"
	"github.com/cjduffett/stork/logger"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

// StorkServer servers the Stork service.
type StorkServer struct {
	Engine  *gin.Engine
	Session *mgo.Session
	Config  *config.StorkConfig
}

// NewServer returns a new StorkServer with the specified StorkConfig.
func NewServer(config *config.StorkConfig) *StorkServer {

	if config.Debug {
		gin.SetMode(gin.DebugMode)
		logger.LogLevel = logger.DebugLevel
	} else {
		gin.SetMode(gin.ReleaseMode)
		logger.LogLevel = logger.DefaultLevel
	}

	return &StorkServer{
		Engine:  gin.Default(),
		Session: nil, // Not instantiated until Run() is called
		Config:  config,
	}
}

// Run starts the StorkServer.
func (s *StorkServer) Run() {
	// Connect to MongoDB
	session, err := mgo.Dial(s.Config.DatabaseHost)
	if err != nil {
		logger.Error("Failed to connect to MongoDB at " + s.Config.DatabaseHost)
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Clone the session to protect the connection
	defer session.Close()
	s.Session = session.Clone()
	logger.Info("Connected to MongoDB at " + s.Config.DatabaseHost)

	// Register middleware (CORS, etc.)
	RegisterMiddleware(s.Engine)

	// Create a new Data Access Layer
	dal := db.NewDataAccessLayer(s.Session, s.Config.DatabaseName)

	// Create a new AWSClient
	awsClient := awsutil.NewAWSClient(s.Config.Debug)

	// Register API routes and setup controllers
	api.RegisterRoutes(s.Engine, dal, awsClient)

	// Start Stork
	logger.Info("Starting Stork on port " + strings.TrimPrefix(s.Config.ServerPort, ":"))
	printStork()

	s.Engine.Run(":" + s.Config.ServerPort)
}

func printStork() {
	logger.Info("Stork version " + config.Version)
	fmt.Println()
	fmt.Println("                     _.--.")
	fmt.Println("                 .-\"`_.--.\\   .-.___________")
	fmt.Println("               .\"_-\"`     \\\\ (  O;------/\\\"'`")
	fmt.Println("             ,.\"=___      =)) \\ \\      /  \\")
	fmt.Println("              `~` .=`~'~)  ( _/ /     /    \\")
	fmt.Println("      =`---====\"\"~`\\          _/     /      \\")
	fmt.Println("                    `-------\"`      /        \\")
	fmt.Println("                                   /          \\")
	fmt.Println("                                  (            )")
	fmt.Println("                                   '._      _.'")
	fmt.Println("                                      '----'")
	fmt.Println()
}
