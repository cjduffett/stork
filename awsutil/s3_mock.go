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
	buckets bucketMap
}

type bucketMap map[string]bool

// NewS3Mock returns a pointer to an initialized S3 mock
func NewS3Mock() *S3Mock {
	return &S3Mock{buckets: make(bucketMap)}
}

// CreateBucket mocks the s3.createBucket operation
func (s *S3Mock) CreateBucket(in *s3.CreateBucketInput) (*s3.CreateBucketOutput, error) {

	// Add it to the list of known buckets if it doesn't already exist
	err := s.addBucket(*in.Bucket)
	if err != nil {
		return nil, errors.New(s3.ErrCodeBucketAlreadyExists)
	}

	// Success response
	return &s3.CreateBucketOutput{
		Location: aws.String("/" + *in.Bucket),
	}, nil
}

// DeleteBucket mocks the s3.deleteBucket operation
func (s *S3Mock) DeleteBucket(in *s3.DeleteBucketInput) (*s3.DeleteBucketOutput, error) {

	// Remove it from the list of known buckets, if it exists
	err := s.removeBucket(*in.Bucket)
	if err != nil {
		return nil, errors.New(s3.ErrCodeNoSuchBucket)
	}
	return &s3.DeleteBucketOutput{}, nil
}

func (s *S3Mock) hasBucket(name string) bool {
	_, ok := s.buckets[name]
	return ok
}

func (s *S3Mock) addBucket(name string) error {
	if s.hasBucket(name) {
		return errors.New("Bucket already exists")
	}
	s.buckets[name] = true
	return nil
}

func (s *S3Mock) removeBucket(name string) error {
	if !s.hasBucket(name) {
		return errors.New("Bucket not found")
	}
	delete(s.buckets, name)
	return nil
}

func (s *S3Mock) clearBuckets() error {
	for bucket := range s.buckets {
		err := s.removeBucket(bucket)
		if err != nil {
			return err
		}
	}
	return nil
}
