package server

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	cors "github.com/itsjamie/gin-cors"
)

// RegisterMiddleware registers all Stork middleware.
func RegisterMiddleware(router *gin.Engine) {
	// CORS middleware
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type, If-Match, If-None-Exist",
		ExposedHeaders:  "Location, ETag, Last-Modified",
		MaxAge:          86400 * time.Second, // Preflight expires after 1 day
		Credentials:     true,
		ValidateHeaders: false,
	}))

	// Redirects
	router.GET("/", RedirectRootToUI)
	router.GET("/stork", RedirectStorkToUI)
}

func RedirectRootToUI(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "/stork/ui")
}

func RedirectStorkToUI(c *gin.Context) {
	c.Redirect(http.StatusPermanentRedirect, "/stork/ui")
}
