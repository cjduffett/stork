package db

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type StateTestSuite struct {
	suite.Suite
}

func TestStateTestSuite(t *testing.T) {
	suite.Run(t, new(StateTestSuite))
}

func (s *StateTestSuite) TestElapsedTime() {
	t := new(Task)
	s.Equal(t.ElapsedTime(), time.Duration(0))

	t.Start()
	s.True(t.ElapsedTime() > time.Duration(0))

	t.End()
	s.Equal(t.ElapsedTime(), t.EndTime.Sub(*t.StartTime))
}
