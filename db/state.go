package db

import "time"

const (
	StatusActive    = "active"
	StatusCompleted = "completed"
	StatusError     = "error"
	StatusAborted   = "aborted"
	StatusDeleted   = "deleted"

	FormatFHIR = "FHIR"
	FormatCCDA = "CCDA"
	FormatHTML = "HTML"
	FormatText = "text"
	FormatCSV  = "CSV"
)

// TaskList is a list of Stork Tasks
type TaskList struct {
	Tasks []Task `json:"tasks"`
}

// Task is a single Stork task
type Task struct {
	ID          string     `bson:"_id" json:"id"`
	Status      string     `bson:"status" json:"status"`
	StartTime   *time.Time `bson:"startTime" json:"startTime"`
	EndTime     *time.Time `bson:"endTime" json:"endTime"`
	InstanceIDs []string   `bson:"instanceIds" json:"instanceIds"`
	BucketName  string     `bson:"bucketName" json:"bucketName"`
	User        string     `bson:"user" json:"user"`
	Formats     []string   `bson:"formats" json:"formats"`
}

// ElapsedTime returns the total runtime for this tasks.
// For active tasks, this changes constantly until the task
// is done or stopped.
func (t *Task) ElapsedTime() time.Duration {
	if t.StartTime == nil {
		// Task has not started yet
		return time.Duration(0)
	}
	if t.EndTime == nil {
		// Task is currently active
		return time.Now().Sub(*t.StartTime)
	}
	// Task is complete
	return t.EndTime.Sub(*t.StartTime)
}

// Start records the current time as the start time
func (t *Task) Start() {
	if t.StartTime == nil {
		now := time.Now()
		t.StartTime = &now
	}
}

// End records the current time as the end time
func (t *Task) End() {
	if t.StartTime != nil && t.EndTime == nil {
		now := time.Now()
		t.EndTime = &now
	}
}
