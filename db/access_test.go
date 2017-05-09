package db

import (
	"testing"

	"gopkg.in/mgo.v2/bson"

	"github.com/cjduffett/stork/config"
	"github.com/cjduffett/stork/logger"
	"github.com/cjduffett/stork/testutil"
	"github.com/stretchr/testify/suite"
	"gopkg.in/mgo.v2"
)

type AccessTestSuite struct {
	testutil.MongoSuite
	session *mgo.Session
	DAL     *StorkDAL
}

func TestAccessTestSuite(t *testing.T) {
	suite.Run(t, new(AccessTestSuite))
}

func (a *AccessTestSuite) SetupSuite() {
	// Verbose logging
	logger.LogLevel = logger.DebugLevel

	// Establish a database session. This session must be closed in
	// TearDownSuite(). The call to DB() has MongoSuite stand up the
	// new database. TearDownDBServer() must also be called in TearDownSuite().
	a.session = a.DB().Session.Copy()

	// Create a new DAL
	config := config.DefaultConfig
	config.DatabaseName = "stork-test"
	a.DAL = NewStorkDAL(a.session, config.DatabaseName)
}

func (a *AccessTestSuite) TearDownTest() {
	// Drop the tasksCollection between tests.
	// This may silently error out if the collection does not exist yet.
	a.DB().C(tasksCollection).DropCollection()
}

func (a *AccessTestSuite) TearDownSuite() {
	// Close the active session
	a.session.Close()

	// Clean up and remove all temporary files from the mocked database.
	// See testutil/mongo_suite.go for more.
	a.TearDownDBServer()
}

func (a *AccessTestSuite) TestCreateTask() {
	var err error

	// Create a new task
	task := &Task{
		ID:          bson.NewObjectId().Hex(),
		Status:      StatusActive,
		InstanceIDs: []string{"abc123", "def456"},
		BucketName:  "test-bucket",
		User:        "bob",
		Formats:     []string{"FHIR", "CSV"},
	}
	task.Start()

	// Try to create it in the database
	taskID, err := a.DAL.CreateTask(task)
	a.NoError(err)
	a.Equal(task.ID, taskID)

	// Check that there is now 1 task in the tasksCollection
	n, err := a.DB().C(tasksCollection).Count()
	a.NoError(err)
	a.Equal(1, n)

	var tasks []Task
	err = a.DB().C(tasksCollection).Find(bson.M{}).All(&tasks)
	a.NoError(err)

	// Make sure all data was preserved during the insert
	createdTask := tasks[0]
	a.Equal(task.ID, createdTask.ID)
	a.Equal(task.Status, createdTask.Status)
	a.Equal(task.InstanceIDs, createdTask.InstanceIDs)
	a.Equal(task.BucketName, createdTask.BucketName)
	a.Equal(task.User, createdTask.User)
	a.Equal(task.Formats, createdTask.Formats)

	// If no ID is present, CreateTask() should assign one
	task.ID = ""
	newTaskID, err := a.DAL.CreateTask(task)
	a.NoError(err)
	a.NotEmpty(newTaskID)
	a.NotEmpty(task.ID)
	a.NotEqual(newTaskID, taskID)

	// There should now be 2 tasks in the database
	n, err = a.DB().C(tasksCollection).Count()
	a.Equal(2, n)
}

func (a *AccessTestSuite) TestGetTask() {
	var err error

	// Create a new task
	task := &Task{
		ID:          bson.NewObjectId().Hex(),
		Status:      StatusActive,
		InstanceIDs: []string{"abc123", "def456"},
		BucketName:  "test-bucket",
		User:        "bob",
		Formats:     []string{"FHIR", "CSV"},
	}
	task.Start()

	// Add it to the database
	taskID, err := a.DAL.CreateTask(task)
	a.NoError(err)

	// Now retrieve it by ID
	gotTask, err := a.DAL.GetTask(taskID)
	a.NoError(err)
	a.NotNil(gotTask)

	// Make sure they're the same
	a.Equal(taskID, gotTask.ID)
	a.Equal(task.Status, gotTask.Status)
}

func (a *AccessTestSuite) TestGetAllTasks() {
	var err error

	// Add some tasks to the database
	task1 := &Task{
		ID:          bson.NewObjectId().Hex(),
		Status:      StatusActive,
		InstanceIDs: []string{"abc123", "def456"},
		BucketName:  "test-bucket-1",
		User:        "bob",
		Formats:     []string{"FHIR", "CSV"},
	}
	task1.Start()

	task1ID, err := a.DAL.CreateTask(task1)
	a.NoError(err)
	a.NotEmpty(task1ID)

	task2 := &Task{
		ID:          bson.NewObjectId().Hex(),
		Status:      StatusActive,
		InstanceIDs: []string{"abc123", "def456"},
		BucketName:  "test-bucket-2",
		User:        "bob",
		Formats:     []string{"FHIR", "CSV"},
	}
	task2.Start()

	task2ID, err := a.DAL.CreateTask(task2)
	a.NoError(err)
	a.NotEmpty(task2ID)

	// Now try to get all tasks
	taskList, err := a.DAL.GetTasks()
	a.NoError(err)
	a.NotNil(taskList)
	a.Len(taskList.Tasks, 2)

	// Now update 1 of the tasks as "deleted"
	task1.Status = StatusDeleted
	_, err = a.DAL.UpdateTask(task1)
	a.NoError(err)

	// Now GetTasks() should return only 1 task
	newTaskList, err := a.DAL.GetTasks()
	a.NoError(err)
	a.NotNil(newTaskList)
	a.Len(newTaskList.Tasks, 1)
}

func (a *AccessTestSuite) TestUpdateTask() {
	var err error

	// Create a new task
	task := &Task{
		ID:          bson.NewObjectId().Hex(),
		Status:      StatusActive,
		InstanceIDs: []string{"abc123", "def456"},
		BucketName:  "test-bucket",
		User:        "bob",
		Formats:     []string{"FHIR", "CSV"},
	}
	task.Start()

	// Add it to the database
	taskID, err := a.DAL.CreateTask(task)
	a.NoError(err)
	a.NotEmpty(taskID)

	// Now update it
	task.User = "geoff"
	task.Status = StatusCompleted
	task.End()

	updatedTask, err := a.DAL.UpdateTask(task)
	a.NoError(err)
	a.NotNil(updatedTask)

	a.Equal(StatusCompleted, updatedTask.Status)
	a.Equal("geoff", updatedTask.User)
	a.NotNil(updatedTask.EndTime)

	// Also get the task from the database and make sure
	// it's really updated there, too
	gotTask, err := a.DAL.GetTask(taskID)
	a.NoError(err)
	a.NotNil(gotTask)

	a.Equal(taskID, gotTask.ID)
	a.Equal(StatusCompleted, gotTask.Status)
	a.Equal("geoff", gotTask.User)
	a.NotNil(gotTask.EndTime)
}

func (a *AccessTestSuite) TestDeleteTask() {

	var err error

	// Create a new task
	task := &Task{
		ID:          bson.NewObjectId().Hex(),
		Status:      StatusActive,
		InstanceIDs: []string{"abc123", "def456"},
		BucketName:  "test-bucket",
		User:        "bob",
		Formats:     []string{"FHIR", "CSV"},
	}
	task.Start()

	// Add it to the database
	taskID, err := a.DAL.CreateTask(task)
	a.NoError(err)
	a.NotEmpty(taskID)

	// Now mark it as deleted
	err = a.DAL.DeleteTask(taskID)
	a.NoError(err)

	// Get the task
	gotTask, err := a.DAL.GetTask(taskID)
	a.NoError(err)
	a.NotNil(gotTask)
	a.Equal(taskID, gotTask.ID)
	a.Equal(StatusDeleted, gotTask.Status)
}
