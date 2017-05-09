package api

import (
	"github.com/cjduffett/stork/config"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

// RegisterRoutes registers all Stork API routes. All API routes are served from /stork/api.
func RegisterRoutes(router *gin.Engine, session *mgo.Session, config *config.StorkConfig) {

	apic := NewAPIController(session, config)

	// All task routes
	taskGroup := router.Group("/task")
	taskGroup.POST("", apic.CreateTask)
	taskGroup.GET("", apic.GetTasks)

	// Specific task item
	taskItem := taskGroup.Group("/:id")
	taskItem.GET("", apic.GetTaskStatus)
	taskItem.DELETE("", apic.DeleteTask)
	taskItem.POST("/abort", apic.AbortTask)

	// Synthea ONLY endpoint
	taskItem.POST("/done", apic.SyntheaInstanceDone)
}
