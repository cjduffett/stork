package site

import (
	"net/http"

	"github.com/cjduffett/stork/config"
	"github.com/cjduffett/stork/logger"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

// RegisterSite registers all Stork UI routes. All site routes are served from /stork/ui
func RegisterSite(router *gin.Engine, session *mgo.Session, config *config.StorkConfig) {
	loadStaticFiles(router, config)
	uiGroup := router.Group("/stork/ui")

	uiGroup.GET("", Index)
}

func loadStaticFiles(router *gin.Engine, config *config.StorkConfig) {
	router.LoadHTMLGlob("assets/*.html")
	router.Static("/static", config.StaticFilePath)
	logger.Debug("Serving static files from ", config.StaticFilePath)
}

// Index serves the UI's main index page.
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
