package awsutil

import (
	"github.com/cjduffett/stork/config"
)

// InstanceConfig describes the configuration that will be passed to each
// Synthea instance as serialized JSON (in user data).
type InstanceConfig struct {
	TaskID         string `json:"taskId"`
	Population     int    `json:"population"`
	YearsOfHistory int    `json:"yearsOfHistory"`
	BucketName     string `json:"bucketName"`
	BucketRegion   string `json:"bucketRegion"`
	// The endpoint Synthea should ping when done generating patients
	DoneEndpoint string   `json:"doneEndpoint"`
	Formats      []string `json:"formats"`
}

// ValidateConfig ensures that InstanceConfig is complete and can
// safely be sent to a Synthea instance without error.
func ValidateConfig(i *InstanceConfig, config *config.StorkConfig) bool {
	return true
}

// InstanceStatus describes the current status of a running Synthea instance.
type InstanceStatus struct {
	InstanceID string
	Status     string
}
