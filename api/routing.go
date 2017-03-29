package api

import (
	"github.com/cjduffett/stork/config"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

// RegisterRoutes registers all Stork API routes. All API routes are served from /stork/api.
func RegisterRoutes(router *gin.Engine, session *mgo.Session, config *config.StorkConfig) {
	apic := NewAPIController(session, config)

	apiGroup := router.Group("/stork/api")
	apiGroup.GET("", apic.APIRoot)
}
