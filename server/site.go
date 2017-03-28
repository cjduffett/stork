package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterSite registers all Stork UI routes. All site routes are served from /stork/ui
func RegisterSite(router *gin.Engine) {
	router.LoadHTMLGlob("assets/*.html")
	router.GET("/stork/ui", Index)
}

// Index serves the UI's main index page.
func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}
