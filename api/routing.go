package api

import (
	"github.com/cjduffett/stork/awsutil"
	"github.com/cjduffett/stork/db"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers all Stork API routes. All API routes are served from /stork/api.
func RegisterRoutes(router *gin.Engine, dal *db.DataAccessLayer, awsClient *awsutil.AWSClient) {

	apic := NewAPIController(dal, awsClient)

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
