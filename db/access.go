package db

import (
	"errors"

	"github.com/cjduffett/stork/logger"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	tasksCollection = "tasks"
)

// StorkDAL is the Stork data access layer, exposing all methods needed
// to access saved state in MongoDB.
type StorkDAL struct {
	session *mgo.Session
	dbname  string
}

// NewStorkDAL creates a new Stork data access layer
func NewStorkDAL(session *mgo.Session, dbname string) *StorkDAL {
	logger.Debug("Creating DAL for database ", dbname)

	return &StorkDAL{
		session: session,
		dbname:  dbname,
	}
}

// GetTask retrieves a Task from the database, by ID
func (s *StorkDAL) GetTask(taskID string) (*Task, error) {
	worker := s.session.Copy()
	defer worker.Close()

	logger.Debug("Getting task ", taskID)

	task := Task{}
	err := worker.DB(s.dbname).C(tasksCollection).FindId(taskID).One(&task)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &task, nil
}

// GetTasks retrieves all tasks from the database, excluding
// those that were deleted.
func (s *StorkDAL) GetTasks() (*TaskList, error) {
	worker := s.session.Copy()
	defer worker.Close()

	logger.Debug("Getting all tasks")

	tasks := []Task{}
	query := bson.M{"status": bson.M{"$not": bson.RegEx{Pattern: StatusDeleted}}}
	err := worker.DB(s.dbname).C(tasksCollection).Find(query).All(&tasks)

	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return &TaskList{Tasks: tasks}, nil
}

// CreateTask adds a new task to the database
func (s *StorkDAL) CreateTask(task *Task) (string, error) {
	worker := s.session.Copy()
	defer worker.Close()

	if task.ID == "" {
		logger.Debug("No task ID found, generating a new one")
		task.ID = bson.NewObjectId().Hex()
	}
	logger.Debug("Creating task ", task.ID)

	err := worker.DB(s.dbname).C(tasksCollection).Insert(*task)

	if err != nil {
		logger.Error(err)
		return "", err
	}
	return task.ID, nil
}

// UpdateTask updates a task in the database
func (s *StorkDAL) UpdateTask(task *Task) (*Task, error) {
	if task.ID == "" {
		// This is an unknown task, error out
		err := errors.New("Unknown task: no task ID found")
		logger.Error(err)
		return nil, err
	}

	worker := s.session.Copy()
	defer worker.Close()

	logger.Debug("Updating task ", task.ID)

	err := worker.DB(s.dbname).C(tasksCollection).UpdateId(task.ID, bson.M{"$set": task})
	if err != nil {
		return nil, err
	}
	return task, nil
}

// DeleteTask marks a task in the database as "deleted"
func (s *StorkDAL) DeleteTask(taskID string) error {
	worker := s.session.Copy()
	defer worker.Close()

	logger.Debug("Marking task ", taskID, " as 'deleted'")
	query := bson.M{"$set": bson.M{"status": StatusDeleted}}
	return worker.DB(s.dbname).C(tasksCollection).UpdateId(taskID, query)
}
