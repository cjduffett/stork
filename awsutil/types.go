package awsutil

// InstanceConfig describes the configuration that will be passed to each
// Synthea instance as serialized JSON (in user data).
type InstanceConfig struct {
	TaskID       string `json:"taskId"`
	Population   int    `json:"population"`
	BucketName   string `json:"bucketName"`
	BucketRegion string `json:"bucketRegion"`
	// The endpoint Synthea should ping when done generating patients
	DoneEndpoint string `json:"doneEndpoint"`
	// The ID of the Synthea AMI to spawn
	ImageID string `json:"imageId"`
}

// InstanceStatus describes the current status of a running Synthea instance.
type InstanceStatus struct{}

// BucketStatus describes the current status of an S3 bucket holding Synthea data.
type BucketStatus struct{}
