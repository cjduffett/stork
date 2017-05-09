package awsutil

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
)

// S3Mock mocks out the AWS S3 API for testing
type S3Mock struct {
	s3iface.S3API
	buckets []string
}

// NewS3Mock returns a pointer to an initialized S3 mock
func NewS3Mock() *S3Mock {
	return &S3Mock{
		buckets: []string{},
	}
}

// CreateBucket mocks the s3.createBucket() operation
func (s *S3Mock) CreateBucket(in *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {

	// Add it to the list of known buckets
	if !s.hasBucket(*in.Bucket) {
		s.buckets = append(s.buckets, *in.Bucket)

		// Success response
		return &s3.CreateBucketOutput{
			Location: aws.String("/" + *in.Bucket),
		}, nil

	}

	// Failure response
	return nil, errors.New(s3.ErrCodeBucketAlreadyExists)
}

func (s *S3Mock) hasBucket(name string) bool {
	for _, bucket := range s.buckets {
		if bucket == name {
			return true
		}
	}
	return false
}
