package awsutil

import (
	"reflect"

	"github.com/cjduffett/stork/config"
)

// InstanceConfig describes the configuration that will be passed to each
// Synthea instance as serialized JSON (in user data).
type InstanceConfig struct {
	TaskID       string `json:"task_id"`
	Population   int    `json:"population"`
	BucketName   string `json:"bucketName"`
	BucketRegion string `json:"bucketRegion"`
	// The endpoint Synthea should ping when done generating patients
	DoneEndpoint string `json:"done_endpoint"`
}

// ValidateConfig ensures that InstanceConfig is complete and can
// safely be sent to a Synthea instance without error.
func ValidateConfig(i *InstanceConfig, config *config.StorkConfig) bool {
	v := reflect.ValueOf(i).Elem() // Use Elem to dereference the pointer

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		switch field.Kind() {
		case reflect.String:
			// The remaining strings must not be empty.
			str := field.String()
			if str == "" {
				return false
			}

		case reflect.Int:
			// This only applies to population:
			// Population must be >= config.MinPopulationSize
			if field.Int() < int64(config.MinPopulationSize) {
				return false
			}

		default:
			// Unknown type in the config object
			return false
		}
	}
	return true
}

// InstanceStatus describes the current status of a running Synthea instance.
type InstanceStatus struct {
	InstanceID string
	Status     string
}
