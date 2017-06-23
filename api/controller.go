package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/cjduffett/stork/awsutil"
	"github.com/cjduffett/stork/db"
	"github.com/cjduffett/stork/logger"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2"
)

// APIController implements all Stork API endpoints
type APIController struct {
	DAL       *db.DataAccessLayer
	AWSClient *awsutil.AWSClient
}

// NewAPIController returns a pointer to an initialized APIController
func NewAPIController(dal *db.DataAccessLayer, awsClient *awsutil.AWSClient) *APIController {
	return &APIController{
		DAL:       dal,
		AWSClient: awsClient,
	}
}

// CreateTask creates a new Stork task, specifying the number of
// patients to generate, number of instances to use, the instance
// type to use, and what formats to export.
func (a *APIController) CreateTask(c *gin.Context) {
	// Parse the request body. The request contains:
	// - population
	// - number of instances
	// - formats to export
	// - bucket to export to
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse("Failed to read request body"))
		return
	}

	req := &CreateTaskRequest{}
	err = json.Unmarshal(body, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, createErrorResponse("Failed to parse request body"))
		return
	}

	// Create a new task ID
	taskID := bson.NewObjectId().Hex()
	logger.Debug("Creating task " + taskID + "...")

	// Create ouput bucket
	err = a.AWSClient.CreateBucket(req.BucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse("Failed to create bucket '"+req.BucketName+"'"))
		return
	}

	// Create EC2 instances
	instanceIDs, err := a.AWSClient.StartInstances(int64(req.NumberOfInstances), &awsutil.InstanceConfig{
		TaskID:         taskID,
		Population:     (req.Population / req.NumberOfInstances),
		BucketName:     req.BucketName,
		BucketRegion:   "us-east-1",
		YearsOfHistory: 5,
		Formats:        req.Formats,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse("Failed to start Synthea instances"))
		return
	}
	c.JSON(http.StatusOK, instanceIDs)

	// Save state
	now := time.Now()
	task := &db.Task{
		ID:          taskID,
		Status:      db.TaskStatusActive,
		StartTime:   &now,
		InstanceIDs: instanceIDs,
		BucketName:  req.BucketName,
		User:        req.User,
		Formats:     req.Formats,
	}
	_, err = a.DAL.CreateTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse("Failed to save task state for task "+taskID))
		return
	}

	// Return status
	c.JSON(http.StatusOK, &CreateTaskResponse{
		TaskID: taskID,
		Status: db.TaskStatusActive,
	})
}

// GetTasks returns a list of all Stork tasks and their statuses.
func (a *APIController) GetTasks(c *gin.Context) {
	// Check state for all non-deleted tasks
	tasks, err := a.DAL.GetTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse("Failed to get tasks: error accessing the database"))
		return
	}

	// return a list of these states & statuses
	c.JSON(http.StatusOK, tasks)
}

// GetTaskStatus gets the current status of an active task. While
// a task is in-progress GetTaskStatus returns the number of running
// instance, the current processing time, etc. Once complete, GetTaskStatus
// returns the URL to the S3 bucket containing all of the exported data.
func (a *APIController) GetTaskStatus(c *gin.Context) {
	// Check state for the desired task
	taskID := c.Param("id")
	task, err := a.DAL.GetTask(taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse("Failed to get status of task '"+taskID+"'"))
		return
	}
	c.JSON(http.StatusOK, task)
}

// AbortTask stops a running Stork task, killing any active instances
// then deleting the S3 bucket used for the export.
func (a *APIController) AbortTask(c *gin.Context) {
	// Check state for active task
	taskID := c.Param("id")
	task, err := a.DAL.GetTask(taskID)
	if err == mgo.ErrNotFound {
		c.JSON(http.StatusNotFound, createErrorResponse("Task '"+taskID+"' not found"))
		return
	}

	if task.Status == db.TaskStatusDeleted {
		c.JSON(http.StatusBadRequest, createErrorResponse("Task already deleted, cannot abort"))
		return
	}

	// Stop running EC2 instances
	err = a.AWSClient.TerminateInstances(task.InstanceIDs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse("Failed to stop instances"))
		return
	}

	// When confirmed stopped, delete the s3 bucket
	err = a.AWSClient.DeleteBucket(task.BucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse("Failed to delete bucket '"+task.BucketName+"'"))
		return
	}

	// Update task status
	task.End()
	updated, err := a.DAL.UpdateTask(task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, createErrorResponse("Failed to delete task '"+taskID+"'"))
		return
	}

	// Return task
	c.JSON(http.StatusOK, updated)
}

// DeleteTask deletes a complete (or aborted) Stork task.
func (a *APIController) DeleteTask(c *gin.Context) {
	// Check state for active task
	taskID := c.Param("id")
	task, err := a.DAL.GetTask(taskID)
	if err == mgo.ErrNotFound {
		c.JSON(http.StatusNotFound, createErrorResponse("Task '"+taskID+"' not found"))
		return
	}

	if task.Status == db.TaskStatusDeleted {
		c.JSON(http.StatusBadRequest, createErrorResponse("Task already deleted"))
		return
	}

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

// ErrorResponse is a response from the API when an error occurs
type ErrorResponse struct {
	Error string `json:"error"`
}

func createErrorResponse(msg string) *ErrorResponse {
	return &ErrorResponse{msg}
}
