package api

import (
	"github.com/cjduffett/stork/config"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

type APIController struct {
	session *mgo.Session
	dbname  string
}

func NewAPIController(session *mgo.Session, config *config.StorkConfig) *APIController {
	return &APIController{
		session: session,
		dbname:  config.DatabaseName,
	}
}

// CreateTask creates a new Stork task, specifying the number of
// patients to generate, number of instances to use, the instance
// type to use, and what formats to export.
func (a *APIController) CreateTask(c *gin.Context) {
	// Read config options:
	// - population
	// - number of instances
	// - instance type (?)
	// - formats to export

	// Create bucket

	// Create EC2 instances

	// Save state

	// Return status
}

// GetTasks returns a list of all Stork tasks and their statuses.
func (a *APIController) GetTasks(c *gin.Context) {
	// Check state for all non-deleted tasks

	// return a list of these states & statuses
}

// GetTaskStatus gets the current status of an active task. While
// a task is in-progress GetTaskStatus returns the number of running
// instance, the current processing time, etc. Once complete, GetTaskStatus
// returns the URL to the S3 bucket containing all of the exported data.
func (a *APIController) GetTaskStatus(c *gin.Context) {
	// Check state for the desired task

	// If not present (or deleted) - error

	// Compute elapsed time

	// Return status
}

// AbortTask stops a running Stork task, killing any active instances
// then deleting the S3 bucket used for the export.
func (a *APIController) AbortTask(c *gin.Context) {
	// Check state for active task

	// If not present (or deleted) - error

	// Stop running EC2 instances

	// When confirmed stopped, delete the s3 bucket

	// Return status
}

// DeleteTask deletes a complete (or aborted) Stork task.
func (a *APIController) DeleteTask(c *gin.Context) {
	// Check state for inactive (or aborted) task

	// If not present (or deleted) - error

	// Delete state and S3 bucket (if applicable)

	// Return confirmation
}

// SyntheaInstanceDone is an endpoint for use by Synthea EC2 instances
// ONLY. Once an instance finishes generating its allocation of patients,
// it pings this endpoint to indicate that it's done.
func (a *APIController) SyntheaInstanceDone(c *gin.Context) {
	// Check state for unique instance ID

	// If not present - error

	// Update state to reflect instance completed, generation count, etc.

	// Return confirmation
}
