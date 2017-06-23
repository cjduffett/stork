package api

// CreateTaskRequest is the POST body of a request to /task.
type CreateTaskRequest struct {
	User              string   `json:"user"`
	Population        int      `json:"population"`
	NumberOfInstances int      `json:"numInstances"`
	Formats           []string `json:"formats"`
	BucketName        string   `json:"bucketName"`
}

// CreateTaskResponse is the body of a response from POST /task.
type CreateTaskResponse struct {
	TaskID string `json:"taskId"`
	Status string `json:"status"`
}

// GetTaskStatusResponse is the body of a response from GET /task/:id
type GetTaskStatusResponse struct {
	TaskID           string `json:"taskId"`
	Status           string `json:"status"`
	ElapsedTime      string `json:"elapsedTime,omitempty"`
	RunningInstances int    `json:"runningInstances,omitempty"`
}
