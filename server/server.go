package server

import (
	"fmt"
	"os"
	"strings"

	"github.com/cjduffett/stork/logger"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

// StorkServer servers the Stork service.
type StorkServer struct {
	Engine  *gin.Engine
	Session *mgo.Session
	Config  StorkConfig
}

// NewServer returns a new StorkServer with the specified StorkConfig.
func NewServer(config StorkConfig) *StorkServer {

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
		logger.Error("Failed to connect to MongoDB")
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Clone the session to protect the connection
	defer session.Close()
	s.Session = session.Clone()
	logger.Info("Connected to MongoDB at " + s.Config.DatabaseHost)

	RegisterMiddleware(s.Engine)
	RegisterRoutes(s.Engine, s.Config)
	RegisterSite(s.Engine)

	// Start Stork
	logger.Info("Starting Stork on port " + strings.TrimPrefix(s.Config.ServerPort, ":"))
	printStork()

	s.Engine.Run(s.Config.ServerPort)
}

func printStork() {
	logger.Info("Stork version " + Version)
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
